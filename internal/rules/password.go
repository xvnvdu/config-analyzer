package rules

import (
	"github.com/xvnvdu/config-analyzer/internal/domain"
)

type RulePassword struct{}

// Проверяем на удовлетворение интерфесу Rule
var _ domain.Rule = (*RulePassword)(nil)

func (r RulePassword) Name() string {
	return "password"
}

func (r RulePassword) Check(cfg domain.Config) []domain.Issue {
	paths := [][]string{
		{"password"},
		{"passwd"},
		{"pswd"},
		{"pass"},
		{"secret"},
		{"private"},
		{"database", "password"},
		{"database", "passwd"},
		{"db", "password"},
	}

	for _, path := range paths {
		val, ok := cfg.Get(path...)
		if !ok {
			continue
		}

		if isOpen(val) {
			return []domain.Issue{
				{
					Level:          domain.MEDIUM,
					Message:        "пароль в открытом виде",
					Recommendation: "Не храните пароль в конфигурационном файле: используйте переменные окружения",
					RuleName:       r.Name(),
				},
			}
		}
	}
	return nil
}

func isOpen(val any) bool {
	str, ok := val.(string)
	return ok && len(str) > 0
}
