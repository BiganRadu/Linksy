package member_service

import (
	internal_models "backend/member_service/models"
	"backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Register handles member registration.
func (m *MemberHandler) Register(c *gin.Context) {
	var member models.Member
	if err := c.BindJSON(&member); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if there is already a member with the same email
	count, err := m.memberDriver.CountMembersWithEmail(member.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(400, gin.H{"error": "A member with this email already exists"})
		return
	}

	// Generate a new ObjectID for the member
	member.ID = primitive.NewObjectID()
	member.CreatedAt = m.nowFunc().Unix()
	token, err := m.tokenHelper.GenerateToken(member.Email, member.Username, member.CreatedAt, time.Hour*24)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	member.Token = token

	// Hash the password before storing it
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

// Login handles member login.
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

	// Generate a new token for the member
	token, err := m.tokenHelper.GenerateToken(foundMember.Email, foundMember.Username, foundMember.CreatedAt, time.Hour*24)
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

// Logout handles member logout by invalidating the token.
func (m *MemberHandler) Logout(c *gin.Context) {
	clientToken := c.Request.Header.Get("AuthToken")
	if clientToken == "" {
		c.JSON(401, gin.H{"error": "Token is required"})
		return
	}
	claims, err := m.tokenHelper.ValidateToken(clientToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}
	userEmail := claims.Email

	// Invalidate the token by setting it to an empty string
	err = m.memberDriver.SetTokenForMember(userEmail, "")
	c.JSON(200, nil)
}

// ChangePassword handles member password change.
func (m *MemberHandler) ChangePassword(c *gin.Context) {
	var changePasswordRequest internal_models.ChangePasswordRequest
	if err := c.BindJSON(&changePasswordRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate the user token
	clientToken := c.Request.Header.Get("AuthToken")
	if clientToken == "" {
		c.JSON(401, gin.H{"error": "Token is required"})
		return
	}
	claims, err := m.tokenHelper.ValidateToken(clientToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	// Get the user email from the token claims
	userEmail := claims.Email
	foundMember, err := m.memberDriver.GetMemberByEmail(userEmail)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check if the old password matches the stored password
	err = bcrypt.CompareHashAndPassword([]byte(foundMember.Password), []byte(changePasswordRequest.OldPassword))
	if err != nil {
		c.JSON(400, gin.H{"error": "Old password is incorrect"})
		return
	}

	// Hash the new password and update the member record
	newPasswordHashed, err := bcrypt.GenerateFromPassword([]byte(changePasswordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	foundMember.Password = string(newPasswordHashed)
	err = m.memberDriver.UpsertMember(foundMember)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Password changed successfully"})
}

// ChangeName handles member name change.
func (m *MemberHandler) ChangeName(c *gin.Context) {
	var changeNameRequest internal_models.ChangeNameRequest
	if err := c.BindJSON(&changeNameRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate the user token
	clientToken := c.Request.Header.Get("AuthToken")
	if clientToken == "" {
		c.JSON(401, gin.H{"error": "Token is required"})
		return
	}
	claims, err := m.tokenHelper.ValidateToken(clientToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	// Get the user email from the token claims
	userEmail := claims.Email
	foundMember, err := m.memberDriver.GetMemberByEmail(userEmail)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Update the member's username
	foundMember.Username = changeNameRequest.NewName
	err = m.memberDriver.UpsertMember(foundMember)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Name changed successfully"})
}
