package testutils

import (
	"database/sql"
	"flag"
	"forum/configs"
	"forum/internal/handlers"
	"forum/internal/render"
	"forum/internal/repository"
	"forum/internal/service"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"os"
	"testing"
)

func NewTestServer(t *testing.T) *handlers.Handler {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred while opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			t.Fatalf("an error '%s' occurred while closing a stub database connection", err)
		}
	}(db)
	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()
	cfg, err := configs.GetConfig(*configPath)
	if err != nil {
		t.Fatalf("error getting config: %v", err)
	}

	// Создаем мок репозитория и сервиса
	mockRepository := repository.NewRepository(db)
	template, err := render.NewTemplateHTML(cfg.TemplateDir)
	if err != nil {
		t.Fatalf("error creating template: %v", err)
	}
	logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	mockService := service.NewService(mockRepository, logger)
	return handlers.NewHandler(mockService, template, cfg.GoogleConfig, cfg.GithubConfig)
}
