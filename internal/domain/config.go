package domain

// Config представляет собой древовидную структуру конфиг-файла
// после анмаршалинга, где корень - Config, ветки - ключи верхнего
// уровня и каждая ветка может быть либо листом (string, bool, т.д.),
// либо узлом (вложенная мапа с новыми ветками)
type Config map[string]any

// Get помогает построить путь от корня Config до листьев
// (переменных конфига). path - это узлы в заданном порядке,
// которые потенциально приведут к листьям
func (c Config) Get(path ...string) (any, bool) {
	var current any = c

	for _, key := range path {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		current, ok = m[key]
		if !ok {
			return nil, false
		}
	}
	return current, true
}
