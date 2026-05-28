package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rabin/tictactoe/internal/api/v1"
	"github.com/rabin/tictactoe/internal/auth"
	"github.com/rabin/tictactoe/internal/storage"
	"github.com/rabin/tictactoe/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting server with config: DB=%s:%s, Port=%s", cfg.Database.Host, cfg.Database.Port, cfg.Server.Port)

	// Initialize database connection
	db, err := storage.NewDatabase(storage.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := storage.RunMigrations(db.DB, "./migrations"); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
		log.Println("Continuing anyway (migrations may already be applied)")
	}

	// Initialize repositories
	userRepo := storage.NewUserRepository(db.DB)
	gameRepo := storage.NewGameRepository(db.DB)

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWT.Secret)

	// Initialize handlers
	authHandler := v1.NewAuthHandler(userRepo, jwtService)
	gameHandler := v1.NewGameHandler(gameRepo, userRepo)
	userHandler := v1.NewUserHandler(userRepo)

	// Create Gin router
	router := gin.Default()

	// Add recovery middleware (catches panics)
	router.Use(gin.Recovery())

	// Add CORS middleware — origins configured via CORS_ALLOWED_ORIGINS env var
	router.Use(cors.New(cors.Config{
		AllowOrigins:  cfg.Server.CORSAllowedOrigins,
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	// Add request logging middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		log.Printf("[%s] %s - %d (%v)", method, path, status, duration)
	})

	// Health check endpoint
	router.GET("/api/v1/health", func(c *gin.Context) {
		if err := db.Health(); err != nil {
			c.JSON(503, gin.H{
				"status":  "error",
				"message": "Database connection failed",
			})
			return
		}
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// API v1 routes - Auth
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.GET("/me", auth.AuthMiddleware(jwtService), authHandler.GetCurrentUser)
	}

	// API v1 routes - Games
	gameRoutes := router.Group("/api/v1/games")
	gameRoutes.Use(auth.AuthMiddleware(jwtService))
	{
		gameRoutes.POST("", gameHandler.CreateGame)
		gameRoutes.GET("/:id", gameHandler.GetGame)
		gameRoutes.POST("/:id/move", gameHandler.MakeMove)
		gameRoutes.GET("", gameHandler.ListGames)
	}

	// API v1 routes - Users
	userRoutes := router.Group("/api/v1/users")
	{
		userRoutes.GET("/:id/stats", userHandler.GetUserStats)
	}

	// Start server in goroutine
	go func() {
		addr := ":" + cfg.Server.Port
		log.Printf("Server listening on %s", addr)
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	fmt.Printf("✓ Server started on port %s\n", cfg.Server.Port)
	fmt.Printf("✓ Health check: http://localhost:%s/api/v1/health\n", cfg.Server.Port)
	fmt.Printf("✓ Auth endpoints:\n")
	fmt.Printf("  - POST http://localhost:%s/api/v1/auth/register\n", cfg.Server.Port)
	fmt.Printf("  - POST http://localhost:%s/api/v1/auth/login\n", cfg.Server.Port)
	fmt.Printf("  - GET  http://localhost:%s/api/v1/auth/me (requires JWT)\n", cfg.Server.Port)
	fmt.Printf("✓ Game endpoints:\n")
	fmt.Printf("  - POST http://localhost:%s/api/v1/games (requires JWT)\n", cfg.Server.Port)
	fmt.Printf("  - GET  http://localhost:%s/api/v1/games/:id (requires JWT)\n", cfg.Server.Port)
	fmt.Printf("  - POST http://localhost:%s/api/v1/games/:id/move (requires JWT)\n", cfg.Server.Port)
	fmt.Printf("  - GET  http://localhost:%s/api/v1/games (requires JWT)\n", cfg.Server.Port)
	fmt.Printf("✓ User endpoints:\n")
	fmt.Printf("  - GET  http://localhost:%s/api/v1/users/:id/stats\n", cfg.Server.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")
	db.Close()
	fmt.Println("Server stopped")
}
