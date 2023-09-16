package middleware

import (
	"bufio"
	"log/slog"
	"os"
	"strings"
)

func LoadEnvVariables() error {
	file, err := os.Open(".env")
	if err != nil {
		slog.Error("Failed to load the environment variable .env", "error", err)
		return err
	}
	defer file.Close()
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}

	return nil
}
