package main

import (
	"backend/app_service"
	"backend/drivers/LinkDriver"
	"backend/drivers/member_driver"
	"backend/helpers"
	"backend/member_service"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	mongoConnectionString := os.Getenv("MONGO_CONNECTION_STRING")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	linkDriver, _ := LinkDriver.NewMongoLinkDriver(mongoConnectionString)
	memberDriver, _ := member_driver.NewMongoMemberDriver(mongoConnectionString)
	tokenHelper := helpers.NewTokenHelper(jwtSecretKey)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow your frontend's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "AuthToken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // Maximum age in seconds
	}))
	appGroup := r.Group("/app")
	appRouter := app_service.NewAppHandler(linkDriver, tokenHelper)
	appRouter.Routes(appGroup)

	memberGroup := r.Group("/member")
	memberRouter := member_service.NewMemberHandler(memberDriver, tokenHelper)
	memberRouter.Routes(memberGroup)

	r.POST("/redirect", appRouter.GetRedirect)

	err := r.Run(":3000")
	if err != nil {
		fmt.Println(err)
	}
}
