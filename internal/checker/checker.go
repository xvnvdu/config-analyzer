package checker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/xvnvdu/config-analyzer/internal/domain"
	"github.com/xvnvdu/config-analyzer/internal/parser"
)

// Checker реализует метод для проверки файлов и хранит
// набор правил, на соответствие которым проверяется файл
type Checker struct {
	rules []domain.Rule
}

// Result представляет собой результат проверки, хранит
// путь к проверямому файлу и слайс проблем с ним
type Result struct {
	Path   string
	Issues []domain.Issue
}

// New создает новый экземпляр Checker
func New(rules []domain.Rule) *Checker {
	return &Checker{rules: rules}
}

// Проверяет файл на соответствие правилам, хранящимся в Checker.
// Может работать как с файлами, так и с директориями, обходя их рекурсивно
func (c *Checker) Check(path string) ([]Result, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		results, err := c.checkDir(path)
		if err != nil {
			return nil, err
		}
		return results, nil
	}

	result, err := c.checkFile(path)
	if err != nil {
		return nil, err
	}
	return []Result{result}, nil
}

// Проверяет cfg на соответствие правилам, хранящимся в Checker
func (c *Checker) CheckConfig(cfg domain.Config) []Result {
	result := Result{Path: "raw"}
	for _, rule := range c.rules {
		result.Issues = append(result.Issues, rule.Check(cfg)...)
	}
	return []Result{result}
}

func (c *Checker) checkFile(path string) (Result, error) {
	result := Result{Path: path}

	result.Issues = append(result.Issues, checkFilePermissions(path)...)

	cfg, err := parser.ParseFile(path)
	if err != nil {
		return result, err
	}

	for _, rule := range c.rules {
		result.Issues = append(result.Issues, rule.Check(cfg)...)
	}
	return result, nil
}

func (c *Checker) checkDir(path string) ([]Result, error) {
	var results []Result
	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(p)
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			return nil
		}
		result, err := c.checkFile(p)
		if err != nil {
			return err
		}
		results = append(results, result)
		return nil
	})
	return results, err
}

// checkFilePermissions проверяет права доступа файла и, при наличии
// проблем, возвращает их в виде ошибки, как остальные правила
func checkFilePermissions(path string) []domain.Issue {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	perm := info.Mode().Perm()
	if perm&0077 != 0 {
		return []domain.Issue{
			{
				Level:          domain.HIGH,
				Message:        fmt.Sprintf("файл конфигурации %s имеет небезопасные права доступа: %o", path, perm),
				Recommendation: fmt.Sprintf("Используйте 'chmod 600 %s' для ограничения прав владельцем", path),
				RuleName:       "file_permissions",
			},
		}
	}
	return nil
}
