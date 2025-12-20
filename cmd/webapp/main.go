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
	"webapp/internal/app/middleware"
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

	// Create temporary storage directory for Docker linters (files are deleted after analysis)
	storageDir := "./storage/tmp"
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("Failed to create temp storage directory: %v", err)
	}

	// Initialize auth
	secretKey := os.Getenv("SESSION_SECRET")
	if secretKey == "" {
		secretKey = "change-this-secret-key-in-production" // Default for development
	}
	middleware.InitAuth(secretKey)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database)
	projectRepo := repositories.NewProjectRepository(database)
	analysisRepo := repositories.NewAnalysisRepository(database)
	issueRepo := repositories.NewIssueRepository(database)
	fileRepo := repositories.NewFileRepository(database)

	// Initialize services
	analyzerService := services.NewAnalyzerService()
	storageService := services.NewStorageService(storageDir)

	// Initialize controllers
	tmpl := parseTemplates()
	authCtrl := controllers.NewAuthController(userRepo, tmpl)
	homeCtrl := controllers.NewHomeController(tmpl)
	uploadCtrl := controllers.NewUploadController(tmpl, analyzerService, storageService, projectRepo, analysisRepo, issueRepo, fileRepo, userRepo)
	analysesCtrl := controllers.NewAnalysesController(tmpl, analysisRepo, projectRepo, issueRepo, fileRepo, userRepo)
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

	// Public routes
	r.HandleFunc("/login", authCtrl.ShowLogin).Methods("GET")
	r.HandleFunc("/login", authCtrl.HandleLogin).Methods("POST")
	r.HandleFunc("/register", authCtrl.ShowRegister).Methods("GET")
	r.HandleFunc("/register", authCtrl.HandleRegister).Methods("POST")
	r.HandleFunc("/logout", authCtrl.HandleLogout).Methods("GET")
	r.HandleFunc("/", homeCtrl.Index).Methods("GET")
	r.HandleFunc("/about", aboutCtrl.Index).Methods("GET")

	// Protected routes
	r.HandleFunc("/upload", middleware.RequireAuth(uploadCtrl.ShowUpload)).Methods("GET")
	r.HandleFunc("/upload/zip", middleware.RequireAuth(uploadCtrl.HandleZipUpload)).Methods("POST")
	r.HandleFunc("/upload/snippet", middleware.RequireAuth(uploadCtrl.HandleSnippetUpload)).Methods("POST")
	r.HandleFunc("/snippet-analyzer", middleware.RequireAuth(uploadCtrl.ShowSnippetAnalyzer)).Methods("GET")
	r.HandleFunc("/api/analyze/snippet", middleware.RequireAuth(uploadCtrl.HandleSnippetAnalyzeAPI)).Methods("POST")
	r.HandleFunc("/analyses", middleware.RequireAuth(analysesCtrl.List)).Methods("GET")
	r.HandleFunc("/analyses/{id}", middleware.RequireAuth(analysesCtrl.Details)).Methods("GET")
	r.HandleFunc("/analyses/{id}/delete", middleware.RequireAuth(analysesCtrl.Delete)).Methods("POST")
	r.HandleFunc("/analyses/{id}/files/{fileId}", middleware.RequireAuth(analysesCtrl.GetFileContent)).Methods("GET")

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

