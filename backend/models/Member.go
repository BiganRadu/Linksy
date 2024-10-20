package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt int64              `bson:"created_at"`
	Token     string             `bson:"token"`
}

type MemberSignedDetails struct {
	Email     string
	Username  string
	CreatedAt int64
	jwt.StandardClaims
}
