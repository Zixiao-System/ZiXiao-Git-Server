package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	// Database drivers
	_ "github.com/lib/pq"               // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3"     // SQLite driver
	_ "github.com/microsoft/go-mssqldb" // SQL Server driver
	"github.com/zixiao/git-server/internal/api"
	"github.com/zixiao/git-server/internal/config"
	"github.com/zixiao/git-server/internal/database"
)

var (
	configPath = flag.String("config", "./configs/server.yaml", "Path to configuration file")
	version    = "1.0.0"
)

func main() {
	flag.Parse()

	// Print banner
	printBanner()

	// Load configuration
	log.Printf("Loading configuration from %s...", *configPath)
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	log.Println("Initializing database...")
	dbCfg := database.Config{
		Type:     cfg.Database.Type,
		Path:     cfg.Database.Path,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		SSLMode:  cfg.Database.SSLMode,
	}
	if err := database.Init(dbCfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create repository directory if it doesn't exist
	if err := os.MkdirAll(cfg.Git.RepoPath, 0755); err != nil {
		log.Fatalf("Failed to create repository directory: %v", err)
	}

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("./logs", 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Setup router
	log.Println("Setting up HTTP router...")
	r := gin.Default()

	// Setup routes
	api.SetupRoutes(r)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting ZiXiao Git Server on %s...", addr)
	log.Printf("Repository path: %s", cfg.Git.RepoPath)
	log.Printf("Database: %s (%s)", cfg.Database.Type, cfg.Database.Path)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════╗
║                                                      ║
║        _____ _  __             _                     ║
║       |__  /(_) \ \  / (_) __ _  ___                ║
║         / / | |  \ \/ /| |/ _' |/ _ \               ║
║        / /_ | |   >  < | | (_| | (_) |              ║
║       /____||_|  /_/\_\|_|\__,_|\___/               ║
║                                                      ║
║              Git Server v%s                      ║
║                                                      ║
╚══════════════════════════════════════════════════════╝
`
	fmt.Printf(banner, version)
	fmt.Println()
}
