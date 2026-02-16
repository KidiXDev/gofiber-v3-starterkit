package utils

import (
	"fmt"
	"os"
)

func ConnectionString() string {
	url := fmt.Sprintf(
		"%s:%s",
		os.Getenv("APP_HOST"),
		os.Getenv("APP_PORT"),
	)

	return url
}

func DatabaseConnectionString() string {
	return os.Getenv("DATABASE_URL")
}
