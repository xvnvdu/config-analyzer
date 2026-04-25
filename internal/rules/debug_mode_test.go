package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xvnvdu/config-analyzer/internal/domain"
)

func TestRuleDebugMode_Check(t *testing.T) {
	issue := &domain.Issue{
		Level:          domain.LOW,
		Message:        "логирование в debug-режиме",
		Recommendation: "Отключите режим отладки или измените его на более избирательный (info+)",
		RuleName:       "debug_mode",
	}

	testCases := []struct {
		name          string
		config        domain.Config
		expectedIssue *domain.Issue
	}{
		{
			name: "Должен найти проблему для log.level: debug",
			config: domain.Config{
				"log": map[string]any{"level": "debug"},
			},
			expectedIssue: issue,
		},
		{
			name: "Должен найти проблему для logging.level: DEBUG (без учета регистра)",
			config: domain.Config{
				"logging": map[string]any{
					"level": "DEBUG",
				},
			},
			expectedIssue: issue,
		},
		{
			name: "Должен найти проблему для debug: true",
			config: domain.Config{
				"debug": true,
			},
			expectedIssue: issue,
		},
		{
			name:          "Не должен находить проблему для безопасного уровня логирования",
			config:        domain.Config{"log": map[string]any{"level": "info"}},
			expectedIssue: nil,
		},
		{
			name:          "Не должен находить проблему для пустого конфига",
			config:        domain.Config{},
			expectedIssue: nil,
		},
		{
			name:          "Не должен находить проблему для несвязанного конфига",
			config:        domain.Config{"server": "prod", "port": 80},
			expectedIssue: nil,
		},
		{
			name: "Не должен находить проблему для debug: false",
			config: domain.Config{
				"debug": false,
			},
			expectedIssue: nil,
		},
	}

	rule := RuleDebugMode{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issues := rule.Check(tc.config)

			if tc.expectedIssue != nil {
				require.NotNil(t, issues, "Ожидалась проблема, но получен nil слайс")
				require.Len(t, issues, 1, "Ожидалась 1 проблема, но получено %d", len(issues))

				require.Equal(t, tc.expectedIssue.Level, issues[0].Level)
				require.Equal(t, tc.expectedIssue.Message, issues[0].Message)
				require.Equal(t, tc.expectedIssue.Recommendation, issues[0].Recommendation)
				require.Equal(t, tc.expectedIssue.RuleName, issues[0].RuleName)
			} else {
				require.Empty(t, issues)
			}
		})
	}
}
