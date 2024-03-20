package main

import (
	"fmt"
	"os"
	"popwa/api"
	"popwa/logs"
	"popwa/repositories"
	"popwa/services"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDatabase()

	popWaRepository := repositories.NewPopWaRepository(db)
	popWaService := services.NewPopWaService(popWaRepository)
	popWahandler := api.NewPopWaHandler(popWaService)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://satjawat.com",
		AllowHeaders: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Get("/getUser/:Username", popWahandler.GetUser)
	app.Post("/addUser", popWahandler.AddUser)
	app.Put("/updateScore", popWahandler.UpdateScore)

	app.Get("/getAllUsers", websocket.New(popWahandler.GetAllUsers))

	app.Listen(fmt.Sprintf(":%v", os.Getenv("POPWA_PORT")))
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logs.Error("Error loading .env file:", zap.Error(err))
	}
}

func initDatabase() *gorm.DB {
	var db *gorm.DB
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	fmt.Println(dsn)
	dial := mysql.Open(dsn)

	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println(db)

	return db
}
