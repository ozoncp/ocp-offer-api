package api

import (
	"context"
	"time"

	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/ozoncp/ocp-offer-api/internal/producer"
	"github.com/ozoncp/ocp-offer-api/internal/repo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api"
)

var (
	totalSuccessCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ocp_offer_api_total_success_created",
		Help: "Total number of requests for offers successfully created",
	})
	totalSuccessUpdated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ocp_offer_api_total_success_updated",
		Help: "Total number of requests for offers successfully updated",
	})
	totalSuccessDeleted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ocp_offer_api_total_success_deleted",
		Help: "Total number of requests for offers successfully deleted",
	})
)

type offerAPI struct {
	pb.UnimplementedOcpOfferApiServiceServer
	repo         repo.IRepository
	dataProducer producer.Producer
}

func NewOfferAPI(r repo.IRepository, p producer.Producer) pb.OcpOfferApiServiceServer {
	return &offerAPI{repo: r, dataProducer: p}
}

func (o *offerAPI) CreateOfferV1(ctx context.Context, req *pb.CreateOfferV1Request) (*pb.CreateOfferV1Response, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreateOfferV1 - invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	offer := models.Offer{
		UserId: req.UserId,
		Grade:  req.Grade,
		TeamId: req.TeamId,
	}

	offerId, err := o.repo.CreateOffer(ctx, offer)
	if err != nil {
		log.Error().Err(err).Msg("CreateOfferV1 -- failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	message := producer.CreateMessage(producer.Create, offerId, time.Now())

	err = o.dataProducer.Send(message)
	if err != nil {
		log.Error().Err(err).Msgf("CreateOfferV1 - failed to send message to kafka")
	}

	totalSuccessCreated.Inc()
	log.Debug().Msgf("CreateOfferV1 - success")

	return &pb.CreateOfferV1Response{
		Id: offerId,
	}, nil
}

// ----------------------------------------------------------------

func (o *offerAPI) MultiCreateOfferV1(ctx context.Context, req *pb.MultiCreateOfferV1Request) (*pb.MultiCreateOfferV1Response, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("MultiCreateOfferV1 - invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	offers := make([]models.Offer, len(req.Offers))

	for i, offer := range offers {
		offers[i] = models.Offer{
			UserId: offer.UserId,
			TeamId: offer.TeamId,
			Grade:  offer.Grade,
		}
	}

	count, err := o.repo.MultiCreateOffer(ctx, offers)
	if err != nil {
		log.Error().Err(err).Msg("MultiCreateOfferV1 -- failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msgf("MultiCreateOfferV1 - success")

	return &pb.MultiCreateOfferV1Response{
		Count: count,
	}, nil
}

// ----------------------------------------------------------------

func (o *offerAPI) DescribeOfferV1(ctx context.Context, req *pb.DescribeOfferV1Request) (*pb.DescribeOfferV1Response, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeOfferV1 - invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	offer, err := o.repo.DescribeOffer(ctx, req.OfferId)
	if err != nil {
		log.Error().Err(err).Msg("DescribeOfferV1 -- failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msgf("DescribeOfferV1 - success")

	return &pb.DescribeOfferV1Response{
		Offer: &pb.Offer{
			Id:     offer.Id,
			UserId: offer.UserId,
			Grade:  offer.Grade,
			TeamId: offer.TeamId,
		},
	}, nil
}

// ----------------------------------------------------------------

func (o *offerAPI) ListOfferV1(ctx context.Context, req *pb.ListOfferV1Request) (*pb.ListOfferV1Response, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("ListOfferV1 - invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	repoOffers, pagInfo, err := o.repo.ListOffer(ctx, models.PaginationInput{
		Cursor: req.Pagination.Cursor,
		Take:   req.Pagination.Take,
		Skip:   req.Pagination.Skip,
	})
	if err != nil {
		log.Error().Err(err).Msg("ListOfferV1 -- failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	offers := make([]*pb.Offer, len(repoOffers))

	for i, val := range repoOffers {
		offers[i] = &pb.Offer{
			Id:     val.Id,
			UserId: val.UserId,
			Grade:  val.Grade,
			TeamId: val.TeamId,
		}
	}

	log.Debug().Msgf("ListOfferV1 - success")

	return &pb.ListOfferV1Response{
		Pagination: &pb.PaginationInfo{
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

// ----------------------------------------------------------------

func (o *offerAPI) UpdateOfferV1(ctx context.Context, req *pb.UpdateOfferV1Request) (*pb.UpdateOfferV1Response, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("UpdateOfferV1 - invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data := models.Offer{
		Id:     req.Id,
		UserId: req.UserId,
		Grade:  req.Grade,
		TeamId: req.TeamId,
	}

	if err := o.repo.UpdateOffer(ctx, data); err != nil {
		log.Error().Err(err).Msg("UpdateOfferV1 -- failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	totalSuccessUpdated.Inc()
	log.Debug().Msgf("UpdateOfferV1 - success")

	return &pb.UpdateOfferV1Response{}, nil
}

func (o *offerAPI) RemoveOfferV1(ctx context.Context, req *pb.RemoveOfferV1Request) (*pb.RemoveOfferV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemoveOfferV1 - invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := o.repo.RemoveOffer(ctx, req.OfferId); err != nil {
		log.Error().Err(err).Msg("RemoveOfferV1 -- failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	totalSuccessDeleted.Inc()
	log.Debug().Msgf("RemoveOfferV1 - success")

	return &pb.RemoveOfferV1Response{}, nil
}
