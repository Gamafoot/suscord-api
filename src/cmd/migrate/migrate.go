package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"suscord/internal/config"
	"suscord/internal/database"
	"suscord/internal/database/gorm/model"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

var tables = []interface{}{
	model.User{},
	model.Session{},
	model.Chat{},
	model.ChatMember{},
	model.Message{},
	model.MessageAttachment{},
	model.Friend{},
}

func main() {
	log.Println("start migrating...")

	cfg := config.GetConfig()

	db, err := database.NewConnect(cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %+v", err)
	}

	if err = db.AutoMigrate(tables...); err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}

	scripts, err := getSqlScripts("function")
	if err != nil {
		log.Fatalf("failed get function scripts: %+v\n", err)
	}

	err = executeScripts(db, scripts)
	if err != nil {
		log.Fatalf("failed to execute functions: %+v\n", err)
	}

	scripts, err = getSqlScripts("trigger")
	if err != nil {
		log.Fatalf("failed get trigger scripts: %+v\n", err)
	}

	err = executeScripts(db, scripts)
	if err != nil {
		log.Fatalf("failed to execute triggers: %+v\n", err)
	}

	log.Println("migrate was successed")
}

func getSqlScripts(folderName string) (map[string]string, error) {
	rootDir := fmt.Sprintf("assets/sql/%s/", folderName)

	scripts := make(map[string]string)

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	for _, entries := range entries {
		if !strings.HasSuffix(entries.Name(), ".sql") {
			continue
		}

		filepath := rootDir + entries.Name()

		content, err := os.ReadFile(filepath)
		if err != nil {
			return nil, pkgErrors.WithStack(err)
		}

		scripts[entries.Name()] = string(content)
	}

	return scripts, nil
}

func executeScripts(db *gorm.DB, scripts map[string]string) error {
	tx := db.Begin()

	for filename, script := range scripts {
		if err := tx.Exec(script).Error; err != nil {
			tx.Rollback()
			return pkgErrors.Errorf("%s: %v", filename, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return pkgErrors.WithStack(err)
	}

	return nil
}
