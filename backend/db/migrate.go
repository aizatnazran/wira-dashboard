package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file from the backend directory
	err := godotenv.Load(filepath.Join("..", ".env"))
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	// Get database connection parameters from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "wira")

	// Connect to the database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Reset the database if requested
	if len(os.Args) > 1 && os.Args[1] == "--reset" {
		log.Println("Resetting database...")
		_, err = db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
		if err != nil {
			log.Fatal("Error resetting database:", err)
		}
		log.Println("Database reset completed!")
	}

	// Get all migration files
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		log.Fatal("Error reading migrations directory:", err)
	}

	// Filter and sort migration files
	var migrations []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, file.Name())
		}
	}
	sort.Strings(migrations)

	// Execute each migration
	for _, migration := range migrations {
		log.Printf("Executing migration: %s", migration)
		
		// Read migration file
		content, err := ioutil.ReadFile(filepath.Join("migrations", migration))
		if err != nil {
			log.Fatalf("Error reading migration file %s: %v", migration, err)
		}

		// Execute the migration
		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("Error executing migration %s: %v", migration, err)
		}

		log.Printf("Successfully executed migration: %s", migration)
	}

	log.Println("All migrations completed successfully")
}


func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}