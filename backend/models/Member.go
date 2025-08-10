package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Member represents a member in the system.
// It contains fields for ID, username, email, password, creation timestamp, and a token.
type Member struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt int64              `bson:"created_at"`
	Token     string             `bson:"token"`
}

// MemberSignedDetails represents the details of a signed member.
// It includes the member's email, username, creation timestamp, and standard JWT claims.
type MemberSignedDetails struct {
	Email     string
	Username  string
	CreatedAt int64
	jwt.StandardClaims
}
