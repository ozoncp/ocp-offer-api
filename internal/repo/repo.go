package repo

import "github.com/ozoncp/ocp-offer-api/internal/models"

type Repo interface {
	AddOffers(offers []models.Offer) error
	ListOffers(limit, offset uint64) ([]models.Offer, error)
	DescribeOffer(offerId uint64) (*models.Offer, error)
}
