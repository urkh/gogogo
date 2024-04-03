package main

import (
	"fmt"
	"log"
	"os"
	"testapp/src/controllers"
	"testapp/src/models"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	db := initDB()
	defer closeDB(db)

	runMigrations(db)
	startServer(db)
}

func startServer(db *gorm.DB) {
	e := echo.New()
	configureMiddleware(e, db)
	registerRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}

func registerRoutes(e *echo.Echo) {
	e.GET("/", controllers.Home)
	e.GET("/show", controllers.Show)

	usersRoute := e.Group("/users")
	usersRoute.POST("/", controllers.NewUser)
	usersRoute.GET("/", controllers.GetUsers)
	usersRoute.GET("/:id/", controllers.GetUser)
	usersRoute.PUT("/:id/", controllers.UpdateUser)
	usersRoute.DELETE("/:id/", controllers.DeleteUser)

	e.Static("/static", "static")
}

func configureMiddleware(e *echo.Echo, db *gorm.DB) {
	allowedHosts := os.Getenv("ALLOWED_HOSTS")

	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{allowedHosts},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}),
		controllers.GormDB(db),
	)
}

func initDB() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	fmt.Println("DB Connected")
	return db
}

func closeDB(db *gorm.DB) {
	connection, _ := db.DB()
	if err := connection.Close(); err != nil {
		log.Fatal("Error closing the database connection:", err)
	}
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
