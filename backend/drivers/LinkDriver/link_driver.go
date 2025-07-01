package LinkDriver

import "backend/models"

type LinkDriver interface {
	GetLinkByID(ID string) (*models.Link, error)
	GetLinksForMember(Email string) ([]*models.Link, error)
	GetQRsForMember(Email string) ([]*models.Link, error)
	UpsertLink(Link *models.Link) error
	DeleteLink(ID string) error
}
