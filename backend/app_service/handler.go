package app_service

import (
	"backend/drivers/LinkDriver"
	"backend/helpers"
	"github.com/gin-gonic/gin"
	"time"
)

type AppHandler struct {
	linkDriver  LinkDriver.LinkDriver
	tokenHelper *helpers.TokenHelper
	nowFunc     func() time.Time
}

func NewAppHandler(linkDriver LinkDriver.LinkDriver, tokenHelper *helpers.TokenHelper) *AppHandler {
	return &AppHandler{
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
	group.POST("/create-link", a.CreateLink)
	group.POST("/update-link", a.UpdateLink)
	group.DELETE("/delete-link", a.DeleteLink)
}
