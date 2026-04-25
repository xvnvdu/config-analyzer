package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xvnvdu/config-analyzer/internal/domain"
)

func TestRuleAlgorithm_Check(t *testing.T) {
	issue := &domain.Issue{
		Level:          domain.HIGH,
		Recommendation: "Замените его на более безопасный",
		RuleName:       "algorithm",
	}

	testCases := []struct {
		name          string
		config        domain.Config
		expectedIssue *domain.Issue
	}{
		{
			name:          "Должен найти проблему для storage.digest-algorithm: SHA-1",
			config:        domain.Config{"storage": map[string]any{"digest-algorithm": "SHA-1"}},
			expectedIssue: issue,
		},
		{
			name:          "Должен найти проблему для crypto.algorithm: RC2",
			config:        domain.Config{"crypto": map[string]any{"algorithm": "RC2"}},
			expectedIssue: issue,
		},
		{
			name:          "Должен найти проблему для security.hash: MD5",
			config:        domain.Config{"security": map[string]any{"hash": "MD5"}},
			expectedIssue: issue,
		},
		{
			name:          "Не должен находить проблему для hash.algorithm: SHA-256",
			config:        domain.Config{"hash": map[string]any{"algorithm": "SHA-256"}},
			expectedIssue: nil,
		},
		{
			name:          "Не должен находить проблему для crypto.algorithm: SHA-3",
			config:        domain.Config{"crypto": map[string]any{"algorithm": "SHA-3"}},
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

	rule := RuleAlgorithm{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issues := rule.Check(tc.config)

			if tc.expectedIssue != nil {
				require.NotNil(t, issues, "Ожидалась проблема, но получен nil слайс")
				require.Len(t, issues, 1, "Ожидалась 1 проблема, но получено %d", len(issues))

				require.Equal(t, tc.expectedIssue.Level, issues[0].Level)
				require.Equal(t, tc.expectedIssue.Recommendation, issues[0].Recommendation)
				require.Equal(t, tc.expectedIssue.RuleName, issues[0].RuleName)
			} else {
				require.Empty(t, issues)
			}
		})
	}
}
