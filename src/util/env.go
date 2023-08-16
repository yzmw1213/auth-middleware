package util

import "os"

func IsDevelopment() bool {
	return os.Getenv("ENV") == "development"
}

func IsProduction() bool {
	return os.Getenv("ENV") == "production"
}
