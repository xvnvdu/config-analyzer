package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/xvnvdu/config-analyzer/internal/domain"
	"gopkg.in/yaml.v3"
)

func ParseFile(path string) (domain.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

func ParseStdin() (domain.Config, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

func Parse(data []byte) (domain.Config, error) {
	var cfg domain.Config

	jsonErr := json.Unmarshal(data, &cfg)
	if jsonErr == nil {
		return cfg, nil
	}
	yamlErr := yaml.Unmarshal(data, &cfg)
	if yamlErr == nil {
		return cfg, nil
	}

	return nil, fmt.Errorf("json: %v, yaml: %v", jsonErr, yamlErr)
}
