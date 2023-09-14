package middleware

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnvVariables() error {
	// Open the .env file.
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file line by line.
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Parse and set the environment variables.
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
