package repo

import (
	"context"
	"fmt"
	"unsafe"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	"github.com/ozoncp/ocp-offer-api/internal/models"
	utils "github.com/ozoncp/ocp-offer-api/internal/utils/models"
)

type IRepository interface {
	MultiCreateOffer(ctx context.Context, offers []models.Offer) (uint64, error)
	CreateOffer(ctx context.Context, offer models.Offer) (uint64, error)
	UpdateOffer(ctx context.Context, offer models.Offer) error
	DescribeOffer(ctx context.Context, offerID uint64) (*models.Offer, error)
	ListOffer(ctx context.Context, pagination models.PaginationInput) ([]models.Offer, *models.PaginationInfo, error)
	RemoveOffer(ctx context.Context, offerID uint64) error
}

type Repository struct {
	db        *sqlx.DB
	batchSize uint
}

func NewRepo(db *sqlx.DB, batchSize uint) IRepository {
	return &Repository{db: db, batchSize: batchSize}
}

func (r *Repository) MultiCreateOffer(ctx context.Context, offers []models.Offer) (uint64, error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("MultiCreateOffer global")
	defer span.Finish()

	var countCreated uint64
	batches, err := utils.SplitOffersToBatches(offers, r.batchSize)
	if err != nil {
		return countCreated, err
	}

	for index, batch := range batches {
		childSpan := tracer.StartSpan(
			fmt.Sprintf("MultiCreateOffer for chunk %d, count of bytes: %d", index, len(batch)*int(unsafe.Sizeof(models.Offer{}))),
			opentracing.ChildOf(span.Context()),
		)
		defer childSpan.Finish()

		query := sq.
			Insert("offer").
			Columns("user_id", "team_id", "grade").
			RunWith(r.db).
			PlaceholderFormat(sq.Dollar)

		for _, offer := range batch {
			query = query.Values(offer.UserID, offer.TeamID, offer.Grade)
		}

		result, err := query.ExecContext(ctx)

		if err != nil {
			return countCreated, err
		}

		rowsAffected, err := result.RowsAffected()

		if err != nil {
			return countCreated, err
		}

		countCreated += uint64(rowsAffected)
	}

	return countCreated, nil
}

func (r *Repository) CreateOffer(ctx context.Context, offer models.Offer) (uint64, error) {
	query := sq.
		Insert("offer").
		Columns("user_id", "team_id", "grade").
		Values(offer.UserID, offer.TeamID, offer.Grade).
		Suffix("RETURNING id").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	var offerID uint64
	if err := query.QueryRowContext(ctx).Scan(&offerID); err != nil {
		return 0, err
	}

	return offerID, nil
}

func (r *Repository) UpdateOffer(ctx context.Context, offer models.Offer) error {
	_, err := sq.
		Update("offer").
		Set("user_id", offer.UserID).
		Set("team_id", offer.TeamID).
		Set("grade", offer.Grade).
		Where(sq.Eq{"id": offer.ID}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)

	return err
}

func (r *Repository) DescribeOffer(ctx context.Context, offerID uint64) (*models.Offer, error) {
	query := sq.
		Select("id", "user_id", "team_id", "grade").
		From("offer").
		Where(sq.And{
			sq.Eq{"id": offerID},
			sq.Eq{"is_deleted": false},
		}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	var offer models.Offer

	err := query.QueryRowContext(ctx).
		Scan(&offer.ID, &offer.UserID, &offer.TeamID, &offer.Grade)

	if err != nil {
		return nil, err
	}

	return &offer, nil
}

func (r *Repository) ListOffer(ctx context.Context, pagination models.PaginationInput) ([]models.Offer, *models.PaginationInfo, error) {
	query := sq.
		Select("id", "user_id", "team_id", "grade").
		From("offer").
		Limit(uint64(pagination.Take)).
		Offset(pagination.Skip).
		OrderBy("id ASC").
		Where(sq.Eq{"is_deleted": false}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	rows, err := query.QueryContext(ctx)

	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	offers := make([]models.Offer, 0)
	for rows.Next() {
		var offer models.Offer
		if err := rows.Scan(
			&offer.ID,
			&offer.UserID,
			&offer.TeamID,
			&offer.Grade,
		); err != nil {
			return nil, nil, err
		}
		offers = append(offers, offer)
	}

	var totalItems uint64
	if err := sq.
		Select("COUNT(*)").
		From("offer").
		Where(sq.Eq{"is_deleted": false}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar).
		QueryRowContext(ctx).Scan(&totalItems); err != nil {
		return nil, nil, err
	}

	perPage := uint32(len(offers))

	pagInfo := pagination.GetPaginationInfo(perPage, totalItems)

	return offers, pagInfo, nil
}

func (r *Repository) RemoveOffer(ctx context.Context, offerID uint64) error {
	_, err := sq.
		Update("offer").
		Set("is_deleted", true).
		Where(sq.Eq{"id": offerID}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)

	return err
}
