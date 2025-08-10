package main

import (
	"backend/app_service"
	"backend/drivers/LinkDriver"
	"backend/drivers/member_driver"
	"backend/helpers"
	"backend/member_service"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	mongoConnectionString := os.Getenv("MONGO_CONNECTION_STRING")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	// Read AWS credentials from environment variables
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_BUCKET_NAME")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)

	if err != nil {
		fmt.Println(err)
	}

	// Instantiate the client with the loaded configuration
	awsClient := s3.NewFromConfig(cfg)
	linkDriver, _ := LinkDriver.NewMongoLinkDriver(mongoConnectionString)
	memberDriver, _ := member_driver.NewMongoMemberDriver(mongoConnectionString)
	tokenHelper := helpers.NewTokenHelper(jwtSecretKey)

	// Initialize the Gin router
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
	// Set up the routes for the app service
	appGroup := r.Group("/app")
	appRouter := app_service.NewAppHandler(awsClient, bucket, linkDriver, tokenHelper)
	appRouter.Routes(appGroup)

	// Set up the routes for the member service
	memberGroup := r.Group("/member")
	memberRouter := member_service.NewMemberHandler(memberDriver, tokenHelper)
	memberRouter.Routes(memberGroup)

	// Set up the redirect route
	// This route is used to handle redirection requests for links
	r.POST("/redirect", appRouter.GetRedirect)

	err = r.Run(":3000")
	if err != nil {
		fmt.Println(err)
	}
}
