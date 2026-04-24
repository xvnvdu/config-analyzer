package domain

type Severity string

const (
	LOW    Severity = "LOW"
	MEDIUM Severity = "MEDIUM"
	HIGH   Severity = "HIGH"
)

// Issue описывает проблему с файлом
// конфигурации, если таковая есть
type Issue struct {
	Level          Severity `json:"level"`
	Message        string   `json:"message"`
	Recommendation string   `json:"recommendation"`
	RuleName       string   `json:"rule_name"`
}
