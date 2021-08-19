package repo

import (
	"github.com/ozoncp/ocp-offer-api/internal/models"
	log "github.com/rs/zerolog/log"
)

type IRepository interface {
	AddOffers(offers []models.Offer) error
	CreateOffer(offer models.Offer) error
	UpdateOffer(offer models.Offer) error
	DescribeOffer(offerId uint64) (*models.Offer, error)
	ListOffer(pagination models.PaginationInput) ([]models.Offer, *models.PaginationInfo, error)
	RemoveOffer(offerId uint64) error
}

type Repository struct{}

func NewRepo() IRepository {
	return &Repository{}
}

func (r *Repository) AddOffers(offers []models.Offer) error {
	log.Printf("AddOffers offers: %v", offers)

	return nil
}

func (r *Repository) CreateOffer(offer models.Offer) error {
	log.Printf("CreateOffer offer: %v", offer)

	return nil
}

func (r *Repository) UpdateOffer(offer models.Offer) error {
	log.Printf("UpdateOffer offer: %v", offer)

	return nil
}

func (r *Repository) DescribeOffer(offerId uint64) (*models.Offer, error) {
	log.Printf("DescribeOffer offerId: %v", offerId)

	return &models.Offer{}, nil
}

func (r *Repository) ListOffer(pagination models.PaginationInput) ([]models.Offer, *models.PaginationInfo, error) {
	log.Printf("ListOffer pagination: %v", pagination)

	pagInfo := &models.PaginationInfo{
		Page:            1,
		TotalPages:      1,
		TotalItems:      1,
		PerPage:         pagination.Take,
		HasNextPage:     false,
		HasPreviousPage: false,
	}

	return []models.Offer{{Id: 1, UserId: 1, Grade: 1, TeamId: 1}}, pagInfo, nil
}

func (r *Repository) RemoveOffer(offerId uint64) error {
	log.Printf("RemoveOffer offerId: %v", offerId)

	return nil
}
