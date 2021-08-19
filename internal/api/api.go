package api

import (
	"context"

	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/ozoncp/ocp-offer-api/internal/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	desc "github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api"
)

type offerAPI struct {
	desc.UnimplementedOcpOfferApiServiceServer
	repo repo.IRepository
}

func NewOfferAPI(repo repo.IRepository) desc.OcpOfferApiServiceServer {
	return &offerAPI{repo: repo}
}

func (o *offerAPI) CreateOfferV1(ctx context.Context, req *desc.CreateOfferV1Request) (*desc.CreateOfferV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	offer := models.Offer{
		UserId: req.UserId,
		Grade:  req.Grade,
		TeamId: req.TeamId,
	}

	if err := o.repo.CreateOffer(offer); err != nil {
		return nil, err
	}

	return &desc.CreateOfferV1Response{}, nil
}

func (o *offerAPI) DescribeOfferV1(ctx context.Context, req *desc.DescribeOfferV1Request) (*desc.DescribeOfferV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	offer, err := o.repo.DescribeOffer(req.OfferId)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &desc.DescribeOfferV1Response{
		Offer: &desc.Offer{
			Id:     offer.Id,
			UserId: offer.UserId,
			Grade:  offer.Grade,
			TeamId: offer.TeamId,
		},
	}, nil
}

func (o *offerAPI) ListOfferV1(ctx context.Context, req *desc.ListOfferV1Request) (*desc.ListOfferV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	repoOffers, pagInfo, err := o.repo.ListOffer(models.PaginationInput{
		Cursor: req.Pagination.Cursor,
		Take:   req.Pagination.Take,
		Skip:   req.Pagination.Skip,
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	offers := make([]*desc.Offer, len(repoOffers))

	for i, val := range repoOffers {
		offers[i] = &desc.Offer{
			Id:     val.Id,
			UserId: val.UserId,
			Grade:  val.Grade,
			TeamId: val.TeamId,
		}
	}

	return &desc.ListOfferV1Response{
		Pagination: &desc.PaginationInfo{
			Page:            pagInfo.Page,
			TotalPages:      pagInfo.TotalPages,
			TotalItems:      pagInfo.TotalItems,
			PerPage:         pagInfo.PerPage,
			HasNextPage:     pagInfo.HasNextPage,
			HasPreviousPage: pagInfo.HasPreviousPage,
		},
		Offers: offers,
	}, nil
}

func (o *offerAPI) UpdateOfferV1(ctx context.Context, req *desc.UpdateOfferV1Request) (*desc.UpdateOfferV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data := models.Offer{
		Id:     req.Id,
		UserId: req.UserId,
		Grade:  req.Grade,
		TeamId: req.TeamId,
	}

	if err := o.repo.UpdateOffer(data); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &desc.UpdateOfferV1Response{}, nil
}

func (o *offerAPI) RemoveOfferV1(ctx context.Context, req *desc.RemoveOfferV1Request) (*desc.RemoveOfferV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := o.repo.RemoveOffer(req.OfferId); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &desc.RemoveOfferV1Response{}, nil
}
