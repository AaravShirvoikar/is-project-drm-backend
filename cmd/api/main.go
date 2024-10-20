package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/handlers"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/repositories"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/services"
	"github.com/AaravShirvoikar/is-project-drm-backend/pkg/auth"
	"github.com/AaravShirvoikar/is-project-drm-backend/pkg/database"
	"github.com/AaravShirvoikar/is-project-drm-backend/pkg/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var (
		serverPort = os.Getenv("PORT")
		dbname     = os.Getenv("DB_DATABASE")
		password   = os.Getenv("DB_PASSWORD")
		username   = os.Getenv("DB_USERNAME")
		dbPort     = os.Getenv("DB_PORT")
		host       = os.Getenv("DB_HOST")
		jwtSecret  = os.Getenv("JWT_SECRET")
		minioHost  = os.Getenv("MINIO_HOST")
		minioPort  = os.Getenv("MINIO_API_PORT")
		accessKey  = os.Getenv("MINIO_ACCESS_KEY")
		secretKey  = os.Getenv("MINIO_SECRET_KEY")
		bucketName = os.Getenv("MINIO_BUCKET_NAME")
	)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, dbPort, dbname)
	db, err := database.NewDatabase(connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	fileStorage, err := storage.NewFileStorage(
		fmt.Sprintf("%s:%s", minioHost, minioPort),
		accessKey,
		secretKey,
		bucketName,
		false,
	)
	if err != nil {
		log.Fatalf("failed to create storage service: %v", err)
	}

	auth.Init(jwtSecret)

	userRepo := repositories.NewUserRepository(db)
	contentRepo := repositories.NewContentRepository(db)
	licenseRepo := repositories.NewLicenseRepository(db)
	sessionKeyRepo := repositories.NewSessionKeyRepo(db)

	userService := services.NewUserService(userRepo)
	contentService := services.NewContentService(contentRepo, fileStorage)
	licenseService := services.NewLicenseService(licenseRepo)
	sessionKeyService := services.NewSessionKeyService(sessionKeyRepo)

	userHandler := handlers.NewUserHandler(userService)
	contentHandler := handlers.NewContentHandler(contentService, licenseService, sessionKeyService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	router.Post("/register", userHandler.Register)
	router.Post("/login", userHandler.Login)

	contentRouter := chi.NewRouter()
	contentRouter.Use(auth.AuthenticateToken)

	contentRouter.Post("/create", contentHandler.CreateContent)

	router.Mount("/content", contentRouter)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", serverPort),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("server started on port %v\n", serverPort)
	log.Fatal(srv.ListenAndServe())
}
