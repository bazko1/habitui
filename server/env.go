package server

import (
	"crypto/rand"
	"log"
	"os"
)

// / nolint:gosec // it is credential name not credential itself.
const secretJWTKeyEnvName = "JWT_SECRET_KEY"

var sqliteDatasePathEnvName = getEnvVariable("SQLITEDB_PATH", "./sqlite.db")

func getEnvVariable(envName, defaultVal string) string {
	if value, ok := os.LookupEnv(envName); ok {
		return value
	}

	return defaultVal
}

func getEnvKey() []byte {
	if secret := getEnvVariable(secretJWTKeyEnvName, ""); secret != "" {
		return []byte(secret)
	}

	randLen := 64
	b := make([]byte, randLen)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("failed to read random bytes", err)
	}

	return b
}
