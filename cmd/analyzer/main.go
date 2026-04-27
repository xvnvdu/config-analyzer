package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	pb "github.com/xvnvdu/config-analyzer/api"
	"github.com/xvnvdu/config-analyzer/internal/checker"
	"github.com/xvnvdu/config-analyzer/internal/domain"
	g "github.com/xvnvdu/config-analyzer/internal/grpc"
	"github.com/xvnvdu/config-analyzer/internal/parser"
	"github.com/xvnvdu/config-analyzer/internal/rules"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Основной набор правил, на соответствие которым
// происходит проверка. Новые добавлять сюда
func defaultRules() []domain.Rule {
	return []domain.Rule{
		rules.RuleAlgorithm{},
		rules.RulePassword{},
		rules.RuleDebugMode{},
		rules.RuleHost{},
		rules.RuleTLS{},
	}
}

func main() {
	help := flag.Bool("h", false, "инструкция по использованию")

	silent := flag.Bool("s", false, "не выходить с ошибкой при наличии проблем")
	flag.BoolVar(silent, "silent", false, "не выходить с ошибкой при наличии проблем")

	stdin := flag.Bool("stdin", false, "читать конфигурацию из stdin")
	mode := flag.String("mode", "cli", "cli | http | grpc")
	port := flag.Int("port", 8080, "порт для http/grpc сервера")

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	rs := defaultRules()
	c := checker.New(rs)

	switch *mode {
	case "cli":
		runCLI(c, *stdin, *silent, flag.Args())
	case "http":
		runHTTP(c, *port)
	case "grpc":
		runGRPC(c, *port)
	default:
		fmt.Fprintln(os.Stderr, "неизвестный режим, используйте флаг '-h' для инструкции по использованию утилиты")
		os.Exit(1)
	}
}

// Запускает утилиту в режиме CLI. Используется по умолчанию, если не был указан другой режим
func runCLI(c *checker.Checker, stdin, silent bool, args []string) {
	if len(args) == 0 && !stdin {
		fmt.Fprintln(os.Stderr, "используйте флаг '-h' для инструкции по использованию утилиты")
		os.Exit(1)
	}

	var result []checker.Result

	if stdin {
		cfg, err := parser.ParseStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка парсинга из stdin: %v\n", err)
			os.Exit(1)
		}
		result = append(result, c.CheckConfig(cfg)...)
	} else {
		path := args[0]

		var err error
		result, err = c.Check(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка парсинга %s: %v\n", path, err)
			os.Exit(1)
		}
	}
	hasIssues := printIssues(result)
	if hasIssues && !silent {
		os.Exit(1)
	}
}

// Запускает утилиту как HTTP-сервер, если при запуске был выбран --mode=http
func runHTTP(c *checker.Checker, port int) {
	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			readErr := fmt.Errorf("ошибка чтения тела запроса: %v", err)
			http.Error(w, readErr.Error(), http.StatusBadRequest)
			return
		}
		cfg, err := parser.Parse(data)
		if err != nil {
			parseErr := fmt.Errorf("ошибка парсинга: %v", err)
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
			return
		}

		result := c.CheckConfig(cfg)
		w.Header().Set("Content-Type", "application/json")

		data, err = json.MarshalIndent(result[0].Issues, "", "  ")
		if err != nil {
			marshalErr := fmt.Errorf("ошибка маршалинга: %v", err)
			http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(data); err != nil {
			fmt.Fprintf(os.Stderr, "ошибка записи: %v\n", err)
		}
	})

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("HTTP сервер запущен на %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка запуска сервера: %v\n", err)
		os.Exit(1)
	}
}

// Запускает утилиту как gRPC-сервер, если при запуске был выбран --mode=grpc
func runGRPC(c *checker.Checker, port int) {
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка запуска listener: %v\n", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAnalyzerServer(grpcServer, g.New(c))
	reflection.Register(grpcServer)

	fmt.Fprintf(os.Stdout, "gRPC сервер запущен на %s\n", addr)
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка запуска сервера: %v\n", err)
		os.Exit(1)
	}
}

// Вспомогательная функция для вывода результатов в режиме CLI,
// возвращает true если была найдена хотя бы одна проблема
func printIssues(result []checker.Result) bool {
	hasIssues := false
	for _, res := range result {
		l := len(res.Issues)
		fmt.Fprintln(os.Stdout, res.Path)
		if l == 0 {
			fmt.Fprintln(os.Stdout, "  └ проблем не найдено")
			continue
		}
		hasIssues = true

		for i, issue := range res.Issues {
			char := "├"
			if i == l-1 {
				char = "└"
			}
			fmt.Fprintf(os.Stdout, "  %s[%s] %s. %s.\n", char, issue.Level, issue.Message, issue.Recommendation)
		}
	}
	return hasIssues
}

// Вспомогательная функция для вывода справки по использованию утилиты
func printHelp() {
	cmd := "analyzer [--mode cli|http|grpc] [--port <порт>] [-s|--silent] [--stdin] [<файл/директория>]\n"
	flags := "Флаги:\n" +
		"  -mode           режим работы утилиты\n" +
		"  -port           порт для запуска утилиты в режиме сервера\n" +
		"  -s, --silent    не выходить с ошибкой при наличии проблем\n" +
		"  -stdin          читать конфигурацию из stdin\n"
	args := "Переменные:\n" +
		"  cli | http | grpc    если указан флаг --mode, по умолчанию 'cli'\n" +
		"  <порт>               если указан флаг --port, по умолчанию '8080'\n" +
		"  <файл/директория>    указывается для проверки конфига только в режиме cli\n"
	exmp := "Примеры:\n" +
		"  'analyzer config.json'                  анализ файла в cli режиме\n" +
		"  'analyzer --silent ./configs/'          анализ директории в cli режиме с флагом -s\n" +
		"  'analyzer --stdin < config.yaml'        анализ файла из стандартного ввода\n" +
		"  'cat config.yaml | analyzer --stdin'    анализ файла из ввода через cat\n" +
		"  'analyzer --mode http --port 9090'      запуск утилиты в режиме http сервера на порту :9090\n" +
		"  'analyzer --mode grpc'                  запуск утилиты в режиме grpc сервера на порту по умолчанию\n"
	info := "Обратите внимание, что флаги/переменные будут работать только в том режиме, для которого они предусмотрены.\n"
	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n%s\n%s", cmd, flags, args, exmp, info)
}
