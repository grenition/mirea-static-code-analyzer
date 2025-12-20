package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/google/uuid"

	"webapp/internal/app/controllers"
	"webapp/internal/app/repositories"
	"webapp/internal/app/services"
	"webapp/internal/db"
)

func main() {
	// Database connection
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/static_analyzer?sslmode=disable"
	}

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Run migrations
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create storage directory
	storageDir := "./storage"
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	// Initialize repositories
	projectRepo := repositories.NewProjectRepository(database)
	analysisRepo := repositories.NewAnalysisRepository(database)
	issueRepo := repositories.NewIssueRepository(database)

	// Initialize services
	analyzerService := services.NewAnalyzerService()
	storageService := services.NewStorageService(storageDir)

	// Initialize controllers
	tmpl := parseTemplates()
	homeCtrl := controllers.NewHomeController(tmpl)
	uploadCtrl := controllers.NewUploadController(tmpl, analyzerService, storageService, projectRepo, analysisRepo, issueRepo)
	analysesCtrl := controllers.NewAnalysesController(tmpl, analysisRepo, projectRepo, issueRepo)
	aboutCtrl := controllers.NewAboutController(tmpl)

	// Setup router
	r := mux.NewRouter()

	// Static files - serve from filesystem
	staticDir := "./web/static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Warning: static directory not found at %s", staticDir)
	} else {
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	}

	// Routes
	r.HandleFunc("/", homeCtrl.Index).Methods("GET")
	r.HandleFunc("/upload", uploadCtrl.ShowUpload).Methods("GET")
	r.HandleFunc("/upload/zip", uploadCtrl.HandleZipUpload).Methods("POST")
	r.HandleFunc("/upload/snippet", uploadCtrl.HandleSnippetUpload).Methods("POST")
	r.HandleFunc("/snippet-analyzer", uploadCtrl.ShowSnippetAnalyzer).Methods("GET")
	r.HandleFunc("/api/analyze/snippet", uploadCtrl.HandleSnippetAnalyzeAPI).Methods("POST")
	r.HandleFunc("/analyses", analysesCtrl.List).Methods("GET")
	r.HandleFunc("/analyses/{id}", analysesCtrl.Details).Methods("GET")
	r.HandleFunc("/about", aboutCtrl.Index).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func parseTemplates() *template.Template {
	tmpl := template.New("").Funcs(template.FuncMap{
		"uuid": func() string {
			return uuid.New().String()
		},
	})

	// Parse all templates together - Go templates will use blocks from all templates
	// We need to parse them in a specific order so layout comes first
	templateDir := "./internal/app/views/templates"
	var templateFiles []string
	
	// Add layout first
	layoutPath := "./internal/app/views/templates/layout.html"
	templateFiles = append(templateFiles, layoutPath)
	
	// Add all other templates
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" && path != layoutPath {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk templates: %v", err)
	}

	tmpl, err = tmpl.ParseFiles(templateFiles...)
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	return tmpl
}

