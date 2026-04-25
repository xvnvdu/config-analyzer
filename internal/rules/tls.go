package rules

import (
	"github.com/xvnvdu/config-analyzer/internal/domain"
)

type RuleTLS struct{}

// Проверяем на удовлетворение интерфесу Rule
var _ domain.Rule = (*RuleTLS)(nil)

func (r RuleTLS) Name() string {
	return "tls"
}

func (r RuleTLS) Check(cfg domain.Config) []domain.Issue {
	paths := [][]string{
		{"tls", "enabled"},
		{"ssl", "enabled"},
		{"tls", "verify"},
		{"ssl", "verify"},
		{"tls", "skip_verify"},
		{"insecure"},
	}

	negativeParams := map[string]struct{}{
		"insecure": {},
		"disabled": {},
	}

	for _, path := range paths {
		val, ok := cfg.Get(path...)
		if !ok {
			continue
		}

		isNegative := false
		for _, p := range path {
			if _, ok := negativeParams[p]; ok {
				isNegative = true
				break
			}
		}

		valIsTrue, ok := val.(bool)
		if !ok {
			continue
		}

		if isNegative && valIsTrue || !isNegative && !valIsTrue {
			return []domain.Issue{
				{
					Level:          domain.HIGH,
					Message:        "TLS отключен",
					Recommendation: "Включите протокол шифрования соединения для повышения безопасности",
					RuleName:       r.Name(),
				},
			}
		}
	}
	return nil
}
