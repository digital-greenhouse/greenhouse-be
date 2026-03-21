package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"digital-greenhouse/greenhouse-be/internal/http/handler"
	"digital-greenhouse/greenhouse-be/internal/http/route"
	"digital-greenhouse/greenhouse-be/internal/repository"
	"digital-greenhouse/greenhouse-be/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	propertyRepo := repository.NewPropertyRepository(db)
	propertySvc := service.NewPropertyService(propertyRepo, userRepo)
	propertyHandler := handler.NewPropertyHandler(propertySvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api/v1", func(r chi.Router) {
		route.RegisterUserRoutes(r, userHandler)
		route.RegisterPropertyRoutes(r, propertyHandler)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor escuchando en http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}

func connectDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "3306"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}

	return db, nil
}
