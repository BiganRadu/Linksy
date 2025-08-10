package member_service

import (
	"backend/drivers/member_driver"
	"backend/helpers"
	"time"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberDriver member_driver.MemberDriver
	tokenHelper  *helpers.TokenHelper
	nowFunc      func() time.Time
}

func NewMemberHandler(memberDriver member_driver.MemberDriver, tokenHelper *helpers.TokenHelper) *MemberHandler {
	return &MemberHandler{
		memberDriver: memberDriver,
		tokenHelper:  tokenHelper,
		nowFunc: func() time.Time {
			return time.Now()
		},
	}
}

// Routes sets up the routes for the member service.
func (m *MemberHandler) Routes(group *gin.RouterGroup) {
	group.POST("/register", m.Register)
	group.POST("change-password", m.ChangePassword)
	group.POST("change-name", m.ChangeName)
	group.POST("/login", m.Login)
	group.GET("/logout", m.Logout)
}
