package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xvnvdu/config-analyzer/internal/domain"
)

func TestRuleTLS_Check(t *testing.T) {
	issue := &domain.Issue{
		Level:          domain.HIGH,
		Message:        "TLS отключен",
		Recommendation: "Включите протокол шифрования соединения для повышения безопасности",
		RuleName:       "tls",
	}

	testCases := []struct {
		name          string
		config        domain.Config
		expectedIssue *domain.Issue
	}{
		{
			name:          "Должен найти проблему для tls.enabled: false",
			config:        domain.Config{"tls": map[string]any{"enabled": false}},
			expectedIssue: issue,
		},
		{
			name:          "Должен найти проблему для ssl.enabled: false",
			config:        domain.Config{"ssl": map[string]any{"enabled": false}},
			expectedIssue: issue,
		},
		{
			name:          "Должен найти проблему для insecure: true",
			config:        domain.Config{"insecure": true},
			expectedIssue: issue,
		},
		{
			name:          "Не должен находить проблему для tls.enabled: true",
			config:        domain.Config{"tls": map[string]any{"enabled": true}},
			expectedIssue: nil,
		},
		{
			name:          "Не должен находить проблему для insecure: false",
			config:        domain.Config{"insecure": false},
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

	rule := RuleTLS{}

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
