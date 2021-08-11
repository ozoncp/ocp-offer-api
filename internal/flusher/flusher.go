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

// Flush добавляет офферы пачками в хранилеще
// при ошибке возвращает не добавленные слайсы и ошибку
func (f *flusher) Flush(offers []models.Offer) ([]models.Offer, error) {
	chunks, err := utils.SplitOffersToBatches(offers, uint(f.chunkSize))
	if err != nil {
		return offers, err
	}

	// позиция - успешно добаленных чанков
	pos := 0

	for _, chunk := range chunks {
		if err := f.offerRepo.AddOffers(chunk); err != nil {
			// возращаем не добавленные в хранилище чанки и ошибку
			return offers[pos:], err
		}
		pos += len(chunk)
	}

	return nil, nil
}
