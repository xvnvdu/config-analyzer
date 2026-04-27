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
	var current any = map[string]any(c)

	// YAML использует вложенные объекты конфига как
	// Config, а не map[string]any, из-за чего некоторые
	// правила падали на !ok, поэтому добавим явное приведение
	// элементов к типу map[string]any для полноты проверок
	for _, key := range path {
		var m map[string]any

		switch v := current.(type) {
		case map[string]any:
			m = v
		case Config:
			m = map[string]any(v)
		default:
			return nil, false
		}

		var ok bool

		current, ok = m[key]
		if !ok {
			return nil, false
		}
	}
	return current, true
}
