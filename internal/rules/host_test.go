package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xvnvdu/config-analyzer/internal/domain"
)

func TestRuleHost_Check(t *testing.T) {
	issue := &domain.Issue{
		Level:          domain.HIGH,
		Message:        "использование 0.0.0.0 открывает сервис на всех сетевых интерфейсах",
		Recommendation: "Не используйте 0.0.0.0 в приватных сервисах без настройки фаервола",
		RuleName:       "host",
	}

	testCases := []struct {
		name          string
		config        domain.Config
		expectedIssue *domain.Issue
	}{
		{
			name:          "Должен найти проблему для host: 0.0.0.0",
			config:        domain.Config{"host": "0.0.0.0"},
			expectedIssue: issue,
		},
		{
			name:          "Должен найти проблему для server.host: 0.0.0.0",
			config:        domain.Config{"server": map[string]any{"host": "0.0.0.0"}},
			expectedIssue: issue,
		},
		{
			name:          "Должен найти проблему для bind: [127.0.0.1, 192.168.0.0, 0.0.0.0]",
			config:        domain.Config{"bind": []any{"127.0.0.1", "192.168.0.0", "0.0.0.0"}},
			expectedIssue: issue,
		},
		{
			name:          "Не должен находить проблему для server.host: [127.0.0.1, 192.168.0.0]",
			config:        domain.Config{"server": map[string]any{"host": []any{"127.0.0.1", "192.168.0.0"}}},
			expectedIssue: nil,
		},
		{
			name:          "Не должен находить проблему для listen: 127.0.0.1",
			config:        domain.Config{"listen": "127.0.0.1"},
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
	}

	rule := RuleHost{}

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
