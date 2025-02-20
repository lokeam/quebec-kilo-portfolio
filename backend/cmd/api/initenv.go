package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Check if file exists AND is not a directory
func envFileExists(filename string) bool {
	fileDetails, err := os.Stat(filename)
	return err == nil && !fileDetails.IsDir()
}

// Attempt to load environment variable file. If override is set, call godotenv.Overload()
func loadEnvFile(filename string, overrideEnvVar bool) {
	if !envFileExists(filename) {
		log.Printf("InitEnv Info: %s not found; skipping loading", filename)
		return
	}

	var envErr error
	if overrideEnvVar {
		envErr = godotenv.Overload(filename)
	} else {
		envErr = godotenv.Load(filename)
	}

	if envErr != nil {
		log.Printf("InitEnv Warn: Unable to load file: %s due to error: %v", filename, envErr)
	} else {
		log.Printf("Successfully loaded file: %s", filename)
	}
}

func initEnv() {
	// Log current working directory
	if workingDir, err := os.Getwd(); err != nil {
		log.Printf("InitEnv Error: unable to identify working directory: %v", err)
	} else {
		log.Printf("Current working directory: %s", workingDir)
	}

	// Get API_ENV flag
	apiEnv := os.Getenv("API_ENV")

	// If we're in dev mode, load local .env files conditionally
	if apiEnv == "dev" {
		loadEnvFile(".env", false) // Load base .env
		loadEnvFile(".env.dev", true) // Override with dev-specific if necessary
	} else {
		log.Println("Prod mode: skipping local .env file. Please ensure env var injection is used.")
	}

	// Validate required environment variables
	requiredVars := []string{
		"API_ENV",
		"DATABASE_URL",
		"PORT",
		"IGDB_CLIENT_ID",
		"IGDB_CLIENT_SECRET",
		"IGDB_AUTH_URL",
	}

	var missingVars []string
	for _, varName := range requiredVars {
		if os.Getenv(varName) == "" {
			missingVars = append(missingVars, varName)
		}
	}

	if len(missingVars) > 0 {
		log.Fatalf("InitEnv Error: Missing required environment variables: %v", missingVars)
	}

	log.Println("InitEnv Info: Environment variables loaded successfully")
}