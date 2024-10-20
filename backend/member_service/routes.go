package member_service

import (
	"backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (m *MemberHandler) Register(c *gin.Context) {
	var member models.Member
	if err := c.BindJSON(&member); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	count, err := m.memberDriver.CountMembersWithEmail(member.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(400, gin.H{"error": "A member with this email already exists"})
		return
	}

	member.ID = primitive.NewObjectID()
	member.CreatedAt = m.nowFunc().Unix()
	token, err := m.tokenHelper.GenerateToken(member.Email, member.Username, member.CreatedAt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	member.Token = token
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(member.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	member.Password = string(hashedPassword)
	err = m.memberDriver.UpsertMember(&member)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, nil)
}

func (m *MemberHandler) Login(c *gin.Context) {
	var member models.Member
	if err := c.BindJSON(&member); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	foundMember, err := m.memberDriver.GetMemberByEmail(member.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if foundMember == nil {
		c.JSON(400, gin.H{"error": "Member not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundMember.Password), []byte(member.Password))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}

	token, err := m.tokenHelper.GenerateToken(foundMember.Email, foundMember.Username, foundMember.CreatedAt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	foundMember.Token = token
	err = m.memberDriver.UpsertMember(foundMember)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": token})
}
