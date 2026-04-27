# Config File Analyzer
**Утилита для анализа конфигурационных файлов веб-приложений и выявления потенциально опасных настроек**

## Функционал
- Парсинг файлов `.json` и `.yaml`/`.yml`
- Проверка файлов по заданному набору правил
- Вывод списка найденных проблем с указанием:
    - Уровня найденной проблемы
    - Описания проблемы
    - Рекомендации по ее устранению

## Набор правил
- Логирование в debug режиме
- Хранение паролей в открытом виде
- Использование `0.0.0.0` без ограничений
- Отключенная TLS-проверка
- Слишком широкие права доступа к файлу
- Использование устаревших/небезопасных алгоритмов

## Возможности
- Работа в режиме CLI:
    - Анализ локальных файлов
    - Рекурсивный анализ файлов в директории
    - Анализ ввода через stdin
- Анализ конфигурации через HTTP API (POST /analyze)
- Анализ конфигурации через gRPC API 
- Справка по использованию утилиты

## Технологии
- **Язык:** Go 1.25+
- **API:** gRPC
- **Веб-серверы:** net/http, grpc

## Установка и запуск
1) Перейдите в директорию, куда хотите сохранить сервис или создайте новую:
```
mkdir <YOUR-DIRECTORY-NAME> && cd <YOUR-DIRECTORY-NAME>
```
2) Клонируйте репозиторий проекта внутри директории и перейдите в него:
```
git clone https://github.com/xvnvdu/config-analyzer.git && cd config-analyzer
```
3) Скомпилируйте бинарь утилиты для удобного использования:
```
go build -o analyzer cmd/analyzer/main.go
```
4) Не знаете, с чего начать ? Используйте флаг `-h`:
```
analyzer -h
```
Вы должны увидеть справку по использованию утилиты:
<details>
<summary>Справка по использованию утилиты</summary>

```
analyzer [--mode cli|http|grpc] [--port <порт>] [-s|--silent] [--stdin] [<файл/директория>]

Флаги:
  -mode           режим работы утилиты
  -port           порт для запуска утилиты в режиме сервера
  -s, --silent    не выходить с ошибкой при наличии проблем
  -stdin          читать конфигурацию из stdin

Переменные:
  cli | http | grpc    если указан флаг --mode, по умолчанию 'cli'
  <порт>               если указан флаг --port, по умолчанию '8080'
  <файл/директория>    указывается для проверки конфига только в режиме cli

Примеры:
  'analyzer config.json'                  анализ файла в cli режиме
  'analyzer --silent ./configs/'          анализ директории в cli режиме с флагом -s
  'analyzer --stdin < config.yaml'        анализ файла из стандартного ввода
  'cat config.yaml | analyzer --stdin'    анализ файла из ввода через cat
  'analyzer --mode http --port 9090'      запуск утилиты в режиме http сервера на порту :9090
  'analyzer --mode grpc'                  запуск утилиты в режиме grpc сервера на порту по умолчанию

Обратите внимание, что флаги/переменные будут работать только в том режиме, для которого они предусмотрены.
```
</details>

## Использование
### CLI-режим

Давайте попробуем проанализировать директорию с конфигами `example-configs`. 
Для этого запустите утилиту с указанием директории в качестве переменной:
```
analyzer example-configs/
```
Вы увидите результаты анализа по каждому файлу в директории:
```
example-configs/config1.json
  ├[HIGH] файл конфигурации example-configs/config1.json имеет небезопасные права доступа: 777. Используйте 'chmod 600 example-configs/config1.json' для ограничения прав владельцем.
  ├[HIGH] слишком слабый алгоритм - SHA1. Замените его на более безопасный.
  ├[MEDIUM] пароль в открытом виде. Не храните пароль в конфигурационном файле: используйте переменные окружения.
  ├[LOW] логирование в debug-режиме. Отключите режим отладки или измените его на более избирательный (info+).
  ├[HIGH] использование 0.0.0.0 открывает сервис на всех сетевых интерфейсах. Не используйте 0.0.0.0 в приватных сервисах без настройки фаервола.
  └[HIGH] TLS отключен. Включите протокол шифрования соединения для повышения безопасности.
example-configs/config1.yaml
  ├[HIGH] слишком слабый алгоритм - MD5. Замените его на более безопасный.
  ├[MEDIUM] пароль в открытом виде. Не храните пароль в конфигурационном файле: используйте переменные окружения.
  ├[LOW] логирование в debug-режиме. Отключите режим отладки или измените его на более избирательный (info+).
  ├[HIGH] использование 0.0.0.0 открывает сервис на всех сетевых интерфейсах. Не используйте 0.0.0.0 в приватных сервисах без настройки фаервола.
  └[HIGH] TLS отключен. Включите протокол шифрования соединения для повышения безопасности.
example-configs/config2.json
  └ проблем не найдено
example-configs/config2.yaml
  └ проблем не найдено
```
Вместо директории вы точно так же можете указать любой локальный файл.

Давайте теперь попробуем прочитать из stdin, для этого используйте соответствующий флаг:
```
cat example-configs/config2.json | analyzer --stdin
```
В этом файле проблем как не было, так и нет:
```
raw
  └ проблем не найдено
```
Если мы попробуем проверить файл с форматом, который не поддерживается утилитой, просто получим ошибку:
```
analyzer cmd/analyzer/main.go
ошибка парсинга cmd/analyzer/main.go: json: invalid character 'p' looking for beginning of value, yaml: yaml: line 55: mapping values are not allowed in this context
```

### HTTP-режим

Теперь запустим утилиту в качестве HTTP-сервера. По умолчанию используется порт `8080`, но давайте выберем другой:
```
analyzer --mode http --port 9090
```
Сервер должен быть успешно запущен, вы увидите уведомление в консоли:
```
HTTP сервер запущен на :9090
```
Утилита ждет, когда мы отправим POST-запрос на эндпоинт `/analyze` с данными для анализа. 
Откройте второй терминал и используйте `curl` для отправки запроса, например, с флагом `--data-binary`, чтобы отправить локальный файл целиком:
```
curl -X POST http://localhost:9090/analyze \
  --data-binary @example-configs/config1.yaml
```
Получим ответ от сервера - список проблем в формате JSON:

<details>
<summary>Ответ сервера</summary>

```json
[
  {
    "level": "HIGH",
    "message": "слишком слабый алгоритм - MD5",
    "recommendation": "Замените его на более безопасный",
    "rule_name": "algorithm"
  },
  {
    "level": "MEDIUM",
    "message": "пароль в открытом виде",
    "recommendation": "Не храните пароль в конфигурационном файле: используйте переменные окружения",
    "rule_name": "password"
  },
  {
    "level": "LOW",
    "message": "логирование в debug-режиме",
    "recommendation": "Отключите режим отладки или измените его на более избирательный (info+)",
    "rule_name": "debug_mode"
  },
  {
    "level": "HIGH",
    "message": "использование 0.0.0.0 открывает сервис на всех сетевых интерфейсах",
    "recommendation": "Не используйте 0.0.0.0 в приватных сервисах без настройки фаервола",
    "rule_name": "host"
  },
  {
    "level": "HIGH",
    "message": "TLS отключен",
    "recommendation": "Включите протокол шифрования соединения для повышения безопасности",
    "rule_name": "tls"
  }
]
```
</details>


### gRPC-режим
Прежде чем запустить gRPC-сервер, установите специальный инструмент для тестирования - `grpcurl`:
```
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

Запустим сервер на порту по умолчанию:
```
analyzer --mode grpc
```

Используем подстановку, чтобы передать файл нашему серверу:
```
CONFIG=$(base64 -w 0 example-configs/config1.json)
grpcurl -plaintext \
  -d "{\"config\": \"$CONFIG\"}" \
  localhost:8080 analyzer.Analyzer/Analyze
```

Снова получим список проблем от сервера:
<details>
<summary>Ответ сервера</summary>

```json
{
  "issues": [
    {
      "level": "HIGH",
      "message": "слишком слабый алгоритм - SHA1",
      "recommendation": "Замените его на более безопасный",
      "ruleName": "algorithm"
    },
    {
      "level": "MEDIUM",
      "message": "пароль в открытом виде",
      "recommendation": "Не храните пароль в конфигурационном файле: используйте переменные окружения",
      "ruleName": "password"
    },
    {
      "level": "LOW",
      "message": "логирование в debug-режиме",
      "recommendation": "Отключите режим отладки или измените его на более избирательный (info+)",
      "ruleName": "debug_mode"
    },
    {
      "level": "HIGH",
      "message": "использование 0.0.0.0 открывает сервис на всех сетевых интерфейсах",
      "recommendation": "Не используйте 0.0.0.0 в приватных сервисах без настройки фаервола",
      "ruleName": "host"
    },
    {
      "level": "HIGH",
      "message": "TLS отключен",
      "recommendation": "Включите протокол шифрования соединения для повышения безопасности",
      "ruleName": "tls"
    }
  ]
}
```
</details>


## Тестирование

Для каждого правила были написаны тесты.

1) Вы можете запустить их следующим образом:
```
go test -coverprofile=coverage.out ./internal/rules/ -v
```
2) А также проверить покрытие в сгенерированном html файле:
```
go tool cover -html=coverage.out
```
Общее покрытие 97.2%:
<details>
<summary>Результаты тестирования правил</summary>

```
go test -coverprofile=coverage.out ./internal/rules/ -v
=== RUN   TestRuleAlgorithm_Check
=== RUN   TestRuleAlgorithm_Check/Должен_найти_проблему_для_storage.digest-algorithm:_SHA-1
=== RUN   TestRuleAlgorithm_Check/Должен_найти_проблему_для_crypto.algorithm:_RC2
=== RUN   TestRuleAlgorithm_Check/Должен_найти_проблему_для_security.hash:_MD5
=== RUN   TestRuleAlgorithm_Check/Не_должен_находить_проблему_для_hash.algorithm:_SHA-256
=== RUN   TestRuleAlgorithm_Check/Не_должен_находить_проблему_для_crypto.algorithm:_SHA-3
=== RUN   TestRuleAlgorithm_Check/Не_должен_находить_проблему_для_пустого_конфига
=== RUN   TestRuleAlgorithm_Check/Не_должен_находить_проблему_для_несвязанного_конфига
--- PASS: TestRuleAlgorithm_Check (0.00s)
    --- PASS: TestRuleAlgorithm_Check/Должен_найти_проблему_для_storage.digest-algorithm:_SHA-1 (0.00s)
    --- PASS: TestRuleAlgorithm_Check/Должен_найти_проблему_для_crypto.algorithm:_RC2 (0.00s)
    --- PASS: TestRuleAlgorithm_Check/Должен_найти_проблему_для_security.hash:_MD5 (0.00s)
    --- PASS: TestRuleAlgorithm_Check/Не_должен_находить_проблему_для_hash.algorithm:_SHA-256 (0.00s)
    --- PASS: TestRuleAlgorithm_Check/Не_должен_находить_проблему_для_crypto.algorithm:_SHA-3 (0.00s)
    --- PASS: TestRuleAlgorithm_Check/Не_должен_находить_проблему_для_пустого_конфига (0.00s)
    --- PASS: TestRuleAlgorithm_Check/Не_должен_находить_проблему_для_несвязанного_конфига (0.00s)
=== RUN   TestRuleDebugMode_Check
=== RUN   TestRuleDebugMode_Check/Должен_найти_проблему_для_log.level:_debug
=== RUN   TestRuleDebugMode_Check/Должен_найти_проблему_для_logging.level:_DEBUG_(без_учета_регистра)
=== RUN   TestRuleDebugMode_Check/Должен_найти_проблему_для_debug:_true
=== RUN   TestRuleDebugMode_Check/Не_должен_находить_проблему_для_безопасного_уровня_логирования
=== RUN   TestRuleDebugMode_Check/Не_должен_находить_проблему_для_пустого_конфига
=== RUN   TestRuleDebugMode_Check/Не_должен_находить_проблему_для_несвязанного_конфига
=== RUN   TestRuleDebugMode_Check/Не_должен_находить_проблему_для_debug:_false
--- PASS: TestRuleDebugMode_Check (0.00s)
    --- PASS: TestRuleDebugMode_Check/Должен_найти_проблему_для_log.level:_debug (0.00s)
    --- PASS: TestRuleDebugMode_Check/Должен_найти_проблему_для_logging.level:_DEBUG_(без_учета_регистра) (0.00s)
    --- PASS: TestRuleDebugMode_Check/Должен_найти_проблему_для_debug:_true (0.00s)
    --- PASS: TestRuleDebugMode_Check/Не_должен_находить_проблему_для_безопасного_уровня_логирования (0.00s)
    --- PASS: TestRuleDebugMode_Check/Не_должен_находить_проблему_для_пустого_конфига (0.00s)
    --- PASS: TestRuleDebugMode_Check/Не_должен_находить_проблему_для_несвязанного_конфига (0.00s)
    --- PASS: TestRuleDebugMode_Check/Не_должен_находить_проблему_для_debug:_false (0.00s)
=== RUN   TestRuleHost_Check
=== RUN   TestRuleHost_Check/Должен_найти_проблему_для_host:_0.0.0.0
=== RUN   TestRuleHost_Check/Должен_найти_проблему_для_server.host:_0.0.0.0
=== RUN   TestRuleHost_Check/Должен_найти_проблему_для_bind:_[127.0.0.1,_192.168.0.0,_0.0.0.0]
=== RUN   TestRuleHost_Check/Не_должен_находить_проблему_для_server.host:_[127.0.0.1,_192.168.0.0]
=== RUN   TestRuleHost_Check/Не_должен_находить_проблему_для_listen:_127.0.0.1
=== RUN   TestRuleHost_Check/Не_должен_находить_проблему_для_пустого_конфига
=== RUN   TestRuleHost_Check/Не_должен_находить_проблему_для_несвязанного_конфига
--- PASS: TestRuleHost_Check (0.00s)
    --- PASS: TestRuleHost_Check/Должен_найти_проблему_для_host:_0.0.0.0 (0.00s)
    --- PASS: TestRuleHost_Check/Должен_найти_проблему_для_server.host:_0.0.0.0 (0.00s)
    --- PASS: TestRuleHost_Check/Должен_найти_проблему_для_bind:_[127.0.0.1,_192.168.0.0,_0.0.0.0] (0.00s)
    --- PASS: TestRuleHost_Check/Не_должен_находить_проблему_для_server.host:_[127.0.0.1,_192.168.0.0] (0.00s)
    --- PASS: TestRuleHost_Check/Не_должен_находить_проблему_для_listen:_127.0.0.1 (0.00s)
    --- PASS: TestRuleHost_Check/Не_должен_находить_проблему_для_пустого_конфига (0.00s)
    --- PASS: TestRuleHost_Check/Не_должен_находить_проблему_для_несвязанного_конфига (0.00s)
=== RUN   TestRulePassword_Check
=== RUN   TestRulePassword_Check/Должен_найти_проблему_для_password:_a1b2c3456
=== RUN   TestRulePassword_Check/Должен_найти_проблему_для_pass:_1234567890
=== RUN   TestRulePassword_Check/Должен_найти_проблему_для_db.password:_qwerty123
=== RUN   TestRulePassword_Check/Не_должен_находить_проблему_для_отключенного_пароля
=== RUN   TestRulePassword_Check/Не_должен_находить_проблему_для_пустого_пароля
=== RUN   TestRulePassword_Check/Не_должен_находить_проблему_для_пустого_конфига
=== RUN   TestRulePassword_Check/Не_должен_находить_проблему_для_несвязанного_конфига
--- PASS: TestRulePassword_Check (0.00s)
    --- PASS: TestRulePassword_Check/Должен_найти_проблему_для_password:_a1b2c3456 (0.00s)
    --- PASS: TestRulePassword_Check/Должен_найти_проблему_для_pass:_1234567890 (0.00s)
    --- PASS: TestRulePassword_Check/Должен_найти_проблему_для_db.password:_qwerty123 (0.00s)
    --- PASS: TestRulePassword_Check/Не_должен_находить_проблему_для_отключенного_пароля (0.00s)
    --- PASS: TestRulePassword_Check/Не_должен_находить_проблему_для_пустого_пароля (0.00s)
    --- PASS: TestRulePassword_Check/Не_должен_находить_проблему_для_пустого_конфига (0.00s)
    --- PASS: TestRulePassword_Check/Не_должен_находить_проблему_для_несвязанного_конфига (0.00s)
=== RUN   TestRuleTLS_Check
=== RUN   TestRuleTLS_Check/Должен_найти_проблему_для_tls.enabled:_false
=== RUN   TestRuleTLS_Check/Должен_найти_проблему_для_ssl.enabled:_false
=== RUN   TestRuleTLS_Check/Должен_найти_проблему_для_insecure:_true
=== RUN   TestRuleTLS_Check/Не_должен_находить_проблему_для_tls.enabled:_true
=== RUN   TestRuleTLS_Check/Не_должен_находить_проблему_для_insecure:_false
=== RUN   TestRuleTLS_Check/Не_должен_находить_проблему_для_пустого_конфига
=== RUN   TestRuleTLS_Check/Не_должен_находить_проблему_для_несвязанного_конфига
--- PASS: TestRuleTLS_Check (0.00s)
    --- PASS: TestRuleTLS_Check/Должен_найти_проблему_для_tls.enabled:_false (0.00s)
    --- PASS: TestRuleTLS_Check/Должен_найти_проблему_для_ssl.enabled:_false (0.00s)
    --- PASS: TestRuleTLS_Check/Должен_найти_проблему_для_insecure:_true (0.00s)
    --- PASS: TestRuleTLS_Check/Не_должен_находить_проблему_для_tls.enabled:_true (0.00s)
    --- PASS: TestRuleTLS_Check/Не_должен_находить_проблему_для_insecure:_false (0.00s)
    --- PASS: TestRuleTLS_Check/Не_должен_находить_проблему_для_пустого_конфига (0.00s)
    --- PASS: TestRuleTLS_Check/Не_должен_находить_проблему_для_несвязанного_конфига (0.00s)
PASS
coverage: 97.2% of statements
ok  	github.com/xvnvdu/config-analyzer/internal/rules	0.003s	coverage: 97.2% of statements
```
</details>

## Полезное

### Архитектура
```
./config-analyzer
├── api               // .proto и сгенерированные pb файлы
├── cmd               // точка входа в приложение
│   └── analyzer
│       └── main.go
├── example-configs   // примеры конфигурационных файлов для проверки
├── go.mod            
├── go.sum            
└── internal
    ├── checker       // модуль проверки файлов - применяет правила из rules
    ├── domain        // основные сущности утилиты - Config, Issue, Rule
    ├── grpc          // grpc сервер
    ├── parser        // парсер файлов/ввода
    └── rules         // основные правила для проверки конфигураций
```
HTTP-сервер было принято реализовать в `main.go` и не выносить его в `internal/http`, 
потому что он содержит всего одну функцию `runHTTP` и один эндпоинт.

С gRPC-сервером ситуация другая - нам нужна отдельная структура, реализующая сгенерированный интерфейс.

### Расширяемость

#### Добавление новых правил
Новые правила добавляются в `internal/rules`, каждое правило реализует пустую структуру:
```go
type RuleMyNewRule struct{}
```
Чтобы применить проверку на новое правило, следует добавить структуру созданного правила в функцию `defaultRules()` внутри `main.go`,
а также реализовать интерфейс `Rule` с методами `Name() string` и `Check(cfg Config) []Issue` в файле нового правила.

Вот и все !

#### Добавление новых путей поиска
Чтобы добавить новые пути поиска по файлу в уже существующее правило, нужно добавить новый слайс путей в переменную `paths`.

Например, мы хотим добавить новый путь для парсинга паролей в базе данных:
```yaml
database:
  settings:
    cache: true
    user: "my_user"
    password: "supersecret777"
```
Для этого в конкретном правиле, проверяющем пароли, мы добавляем новый путь:
```go
paths := [][]string{
    // ... все остальные пути
    {"database", "settings", "password"},
}
```

#### Добавление новых форматов
Проверка формата реализована в функции `Parse` внутри парсера `internal/parser/parser.go`.

Для добавления нового формата достаточно отредактировать эту функцию.

### Проверка прав доступа
Проверка прав реализована отдельно от основных правил и находится в `internal/checker/checker.go`, 
грубо говоря, являясь правилом "по умолчанию". 
Это правило работает только в режиме CLI для проверки отдельных файлов и директорий рекурсивно.

Для получения прав используется `os.Stat`

Проверка производится с использованием маски `0077`, исходя из той логики, что любой доступ для кого-либо, кроме владельца, - это плохо.

То есть, мы проверяем файл на наличие ХОТЯ БЫ ОДНОГО права у группы или остальных, и, если оно есть, отмечаем как опасность, например:
```
  110 110 110   (права 0666) -rw-rw-rw-
& 000 111 111   (маска 0077)
-------------
= 000 110 110   (Результат 0066. Если не 0000 - опасность)
```
