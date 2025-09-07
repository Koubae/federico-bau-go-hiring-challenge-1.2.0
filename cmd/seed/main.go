package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/app/models"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mytheresa/go-hiring-challenge/app/database"
)

const DryRun = false

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db, err := database.New(false)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Shutdown()

	session := db.DB.Session(
		&gorm.Session{
			DryRun:      DryRun,
			Logger:      logger.Default.LogMode(db.Config.LogLevel),
			PrepareStmt: false,
		},
	)

	dir := os.Getenv("POSTGRES_SQL_DIR")
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("reading directory failed: %v", err)
	}

	// DROP ALL TABLES
	log.Println("Dropping all tables")
	if err := session.Migrator().DropTable(&models.Category{}, &models.Product{}, &models.Variant{}); err != nil {
		panic(err)
	}
	log.Println("All database tables dropped successfully!")

	log.Println("Creating all tables")
	if err := session.AutoMigrate(&models.Category{}, &models.Product{}, &models.Variant{}); err != nil {
		panic(err)
	}
	log.Println("All database tables created successfully!")

	log.Println("Seeding database with test data")
	// Filter and sort .sql files
	var sqlFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file)
		}
	}
	sort.Slice(
		sqlFiles, func(i, j int) bool {
			return sqlFiles[i].Name() < sqlFiles[j].Name()
		},
	)

	for _, file := range sqlFiles {
		path := filepath.Join(dir, file.Name())

		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("reading file %s failed: %v", file.Name(), err)
		}

		sql := string(content)
		if err := session.Exec(sql).Error; err != nil {
			log.Printf("executing %s failed: %v", file.Name(), err)
			return
		}

		log.Printf("Executed %s successfully\n", file.Name())
	}
}
