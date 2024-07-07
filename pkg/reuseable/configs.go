package reuseable

import (
	"os"
)

// ///////////////////////////////////////////
// Load env file and get value
// ///////////////////////////////////////////
func GetEnvVar(theValue string) string {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	return os.Getenv(theValue)
}
