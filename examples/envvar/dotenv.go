package main

import (
	"fmt"
	"os"
)

type EnvVarParameterStore struct{}

func (ps EnvVarParameterStore) GetParams(paramMap map[string]string) error {
	for key := range paramMap {
		value := os.Getenv(key)
		if value == "" {
			return fmt.Errorf(`environment variable %s does not exist`, key)
		}
		paramMap[key] = value
	}
	return nil
}
