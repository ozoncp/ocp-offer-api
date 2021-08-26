package flusher

import (
	"context"

	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/ozoncp/ocp-offer-api/internal/repo"
	utils "github.com/ozoncp/ocp-offer-api/internal/utils/models"
)

// Flusher - интерфейс для сброса задач в хранилище.
type Flusher interface {
	Flush(ctx context.Context, offers []models.Offer) ([]models.Offer, error)
}

// NewFlusher возвращает Flusher с поддержкой батчевого сохранения.
func NewFlusher(chunkSize int, offerRepo repo.IRepository) Flusher {
	return &flusher{
		chunkSize: chunkSize,
		offerRepo: offerRepo,
	}
}

type flusher struct {
	chunkSize int
	offerRepo repo.IRepository
}

// Flush добавляет офферы пачками в хранилеще
// при ошибке возвращает не добавленные слайсы и ошибку.
func (f *flusher) Flush(ctx context.Context, offers []models.Offer) ([]models.Offer, error) {
	chunks, err := utils.SplitOffersToBatches(offers, uint(f.chunkSize))
	if err != nil {
		return offers, err
	}

	// позиция - успешно добаленных чанков
	pos := 0

	for _, chunk := range chunks {
		if _, err := f.offerRepo.MultiCreateOffer(ctx, chunk); err != nil {
			// возращаем не добавленные в хранилище чанки и ошибку
			return offers[pos:], err
		}
		pos += len(chunk)
	}

	return nil, nil
}
