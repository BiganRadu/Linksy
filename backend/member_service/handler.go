package member_service

import (
	"backend/drivers/member_driver"
	"backend/helpers"
	"github.com/gin-gonic/gin"
	"time"
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

func (m *MemberHandler) Routes(group *gin.RouterGroup) {
	group.POST("/register", m.Register)
	group.POST("/login", m.Login)
}
