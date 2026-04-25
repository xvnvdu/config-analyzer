package rules

import (
	"fmt"
	"strings"

	"github.com/xvnvdu/config-analyzer/internal/domain"
)

type RuleAlgorithm struct{}

// Проверяем на удовлетворение интерфесу Rule
var _ domain.Rule = (*RuleAlgorithm)(nil)

func (r RuleAlgorithm) Name() string {
	return "algorithm"
}

func (r RuleAlgorithm) Check(cfg domain.Config) []domain.Issue {
	paths := [][]string{
		{"storage", "digest-algorithm"},
		{"crypto", "algorithm"},
		{"hash", "algorithm"},
		{"security", "hash"},
	}

	for _, path := range paths {
		val, ok := cfg.Get(path...)
		if !ok {
			continue
		}

		if isWeakAlgorithm(val) {
			return []domain.Issue{
				{
					Level:          domain.HIGH,
					Message:        fmt.Sprintf("слишком слабый алгоритм - %s", val),
					Recommendation: "Замените его на более безопасный",
					RuleName:       r.Name(),
				},
			}
		}
	}
	return nil
}

func isWeakAlgorithm(val any) bool {
	weakAlgorithms := map[string]struct{}{
		"md5":   {},
		"sha1":  {},
		"sha-1": {},
		"des":   {},
		"rc2":   {},
		"rc4":   {},
	}

	algo, ok := val.(string)
	if ok {
		if _, ok := weakAlgorithms[strings.ToLower(algo)]; ok {
			return true
		}
	}
	return false
}
