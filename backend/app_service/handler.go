package app_service

import (
	"backend/drivers/LinkDriver"
	"backend/helpers"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"time"
)

type AppHandler struct {
	awsClient   *s3.Client
	s3Bucket    string
	linkDriver  LinkDriver.LinkDriver
	tokenHelper *helpers.TokenHelper
	nowFunc     func() time.Time
}

func NewAppHandler(awsClient *s3.Client, s3Bucket string, linkDriver LinkDriver.LinkDriver, tokenHelper *helpers.TokenHelper) *AppHandler {
	return &AppHandler{
		awsClient:   awsClient,
		s3Bucket:    s3Bucket,
		linkDriver:  linkDriver,
		tokenHelper: tokenHelper,
		nowFunc: func() time.Time {
			return time.Now()
		},
	}
}

func (a *AppHandler) Routes(group *gin.RouterGroup) {
	group.Use(a.AuthenticationMiddleware)
	group.GET("/member-info", a.GetMemberInfo)
	group.GET("/link", a.GetLink)
	group.GET("/member-links", a.GetMemberLinks)
	group.GET("/member-qrs", a.GetMemberQRs)
	group.POST("/create-link", a.CreateLink)
	group.POST("/update-link", a.UpdateLink)
	group.DELETE("/delete-link", a.DeleteLink)
	group.DELETE("/delete-qr", a.DeleteQr)
}
