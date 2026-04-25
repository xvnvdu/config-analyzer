package rules

import (
	"strings"

	"github.com/xvnvdu/config-analyzer/internal/domain"
)

type RuleDebugMode struct{}

// Проверяем на удовлетворение интерфесу Rule
var _ domain.Rule = (*RuleDebugMode)(nil)

func (r RuleDebugMode) Name() string {
	return "debug_mode"
}

func (r RuleDebugMode) Check(cfg domain.Config) []domain.Issue {
	paths := [][]string{
		{"log", "level"},
		{"logs", "level"},
		{"logging", "level"},
		{"debug"}, // если debug: true
	}

	for _, path := range paths {
		val, ok := cfg.Get(path...)
		if !ok {
			continue
		}

		if isDebug(val) {
			return []domain.Issue{
				{
					Level:          domain.LOW,
					Message:        "логирование в debug-режиме",
					Recommendation: "Отключите режим отладки или измените его на более избирательный (info+)",
					RuleName:       r.Name(),
				},
			}
		}
	}
	return nil
}

func isDebug(val any) bool {
	switch v := val.(type) {
	case string:
		return strings.EqualFold(v, "debug")
	case bool:
		return v
	default:
		return false
	}
}
