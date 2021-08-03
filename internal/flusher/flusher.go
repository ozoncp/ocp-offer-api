package flusher

import (
	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/ozoncp/ocp-offer-api/internal/repo"
	utils "github.com/ozoncp/ocp-offer-api/internal/utils/models"
)

// Flusher - интерфейс для сброса задач в хранилище
type Flusher interface {
	Flush(offers []models.Offer) ([]models.Offer, error)
}

// NewFlusher возвращает Flusher с поддержкой батчевого сохранения
func NewFlusher(chunkSize int, offerRepo repo.Repo) Flusher {
	return &flusher{
		chunkSize: chunkSize,
		offerRepo: offerRepo,
	}
}

type flusher struct {
	chunkSize int
	offerRepo repo.Repo
}

func (f *flusher) Flush(offers []models.Offer) ([]models.Offer, error) {
	batches, err := utils.SplitOffersToBatches(offers, uint(f.chunkSize))
	if err != nil {
		return nil, err
	}

	for _, chunk := range batches {
		if err := f.offerRepo.AddOffers(chunk); err != nil {
			return nil, err
		}
	}

	return offers, nil
}
