package rules

import (
	"github.com/xvnvdu/config-analyzer/internal/domain"
)

type RuleHost struct{}

// Проверяем на удовлетворение интерфесу Rule
var _ domain.Rule = (*RuleHost)(nil)

func (r RuleHost) Name() string {
	return "host"
}

func (r RuleHost) Check(cfg domain.Config) []domain.Issue {
	paths := [][]string{
		{"host"},
		{"server", "host"},
		{"service", "host"},
		{"listen"},
		{"bind"},
		{"address"},
	}

	for _, path := range paths {
		val, ok := cfg.Get(path...)
		if !ok {
			continue
		}

		if containsUnsafeHost(val) {
			return []domain.Issue{
				{
					Level:          domain.HIGH,
					Message:        "использование 0.0.0.0 открывает сервис на всех сетевых интерфейсах",
					Recommendation: "Не используйте 0.0.0.0 в приватных сервисах без настройки фаервола",
					RuleName:       r.Name(),
				},
			}
		}
	}
	return nil
}

func containsUnsafeHost(val any) bool {
	switch v := val.(type) {
	case string:
		return v == "0.0.0.0"
	case []any:
		for _, host := range v {
			if str, ok := host.(string); ok && str == "0.0.0.0" {
				return true
			}
		}
	}
	return false
}
