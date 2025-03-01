package member_driver

import "backend/models"

type MemberDriver interface {
	CountMembersWithEmail(Email string) (int, error)
	GetMemberByEmail(Email string) (*models.Member, error)
	UpsertMember(Member *models.Member) error
	SetTokenForMember(Email, Token string) error
	DeleteMember(Email string) error
}
