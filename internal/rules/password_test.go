package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xvnvdu/config-analyzer/internal/domain"
)

func TestRulePassword_Check(t *testing.T) {
	issue := &domain.Issue{
		Level:          domain.MEDIUM,
		Message:        "пароль в открытом виде",
		Recommendation: "Не храните пароль в конфигурационном файле: используйте переменные окружения",
		RuleName:       "password",
	}

	testCases := []struct {
		name          string
		config        domain.Config
		expectedIssue *domain.Issue
	}{
		{
			name:          "Должен найти проблему для password: a1b2c3456",
			config:        domain.Config{"password": "a1b2c3456"},
			expectedIssue: issue,
		},
		{
			name:          "Должен найти проблему для pass: 1234567890",
			config:        domain.Config{"pass": "1234567890"},
			expectedIssue: issue,
		},
		{
			name:          "Должен найти проблему для db.password: qwerty123",
			config:        domain.Config{"db": map[string]any{"password": "qwerty123"}},
			expectedIssue: issue,
		},
		{
			name:          "Не должен находить проблему для отключенного пароля",
			config:        domain.Config{"password": false},
			expectedIssue: nil,
		},
		{
			name:          "Не должен находить проблему для пустого пароля",
			config:        domain.Config{"password": ""},
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

	rule := RulePassword{}

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
