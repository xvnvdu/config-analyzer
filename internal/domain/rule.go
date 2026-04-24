package domain

// Rule определяет поведение каждого правила для проверки конфиг-файлов
type Rule interface {
	Name() string
	Check(cfg Config) []Issue
}
