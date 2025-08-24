package env

import (
	"os"
	"strconv"
)

func GetEnvString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func GetEnvInt(key string, defaultValue int) int{
	if value, exists := os.LookupEnv(key); exists{
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue 
		}
	}
	return defaultValue
}