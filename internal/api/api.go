package api

import (
	"context"

	desc "github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api"
	log "github.com/rs/zerolog/log"
)

type offerAPI struct {
	desc.UnimplementedOcpOfferApiServiceServer
}

func NewOfferAPI() desc.OcpOfferApiServiceServer {
	return &offerAPI{}
}

func (o *offerAPI) CreateOfferV1(ctx context.Context, req *desc.CreateOfferV1Request) (*desc.CreateOfferV1Response, error) {
	log.Printf("Create offer: %v", req)

	if err := req.Validate(); err != nil {
		log.Debug().Err(err)
		return nil, err
	}

	return &desc.CreateOfferV1Response{}, nil
}

func (o *offerAPI) DescribeOfferV1(ctx context.Context, req *desc.DescribeOfferV1Request) (*desc.DescribeOfferV1Response, error) {
	log.Printf("Describe offer: %v", req)

	if err := req.Validate(); err != nil {
		log.Debug().Err(err)
		return nil, err
	}

	return &desc.DescribeOfferV1Response{
		Offer: &desc.Offer{
			Id:     1,
			UserId: 2,
			Grade:  3,
			TeamId: 4,
		},
	}, nil
}

func (o *offerAPI) ListOfferV1(ctx context.Context, req *desc.ListOfferV1Request) (*desc.ListOfferV1Response, error) {
	log.Printf("List offers: %v", req)

	if err := req.Validate(); err != nil {
		log.Debug().Err(err)
		return nil, err
	}

	return &desc.ListOfferV1Response{
		Pagination: &desc.PaginationInfo{
			Page:            1,
			TotalPages:      1,
			TotalItems:      0,
			PerPage:         0,
			HasNextPage:     false,
			HasPreviousPage: false,
		},
		Offers: make([]*desc.Offer, 0),
	}, nil
}

func (o *offerAPI) UpdateOfferV1(ctx context.Context, req *desc.UpdateOfferV1Request) (*desc.UpdateOfferV1Response, error) {
	log.Printf("Update offer: %v", req)

	if err := req.Validate(); err != nil {
		log.Debug().Err(err)
		return nil, err
	}

	return &desc.UpdateOfferV1Response{}, nil
}

func (o *offerAPI) RemoveOfferV1(ctx context.Context, req *desc.RemoveOfferV1Request) (*desc.RemoveOfferV1Response, error) {
	log.Printf("Remove offer: %v", req)

	if err := req.Validate(); err != nil {
		log.Debug().Err(err)
		return nil, err
	}

	return &desc.RemoveOfferV1Response{}, nil
}
