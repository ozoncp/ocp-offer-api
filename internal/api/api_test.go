package api_test

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozoncp/ocp-offer-api/internal/api"
	"github.com/ozoncp/ocp-offer-api/internal/mocks"
	"github.com/ozoncp/ocp-offer-api/internal/models"
	pb "github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

var _ = Describe("OcpOfferApiService", func() {

	var (
		listener  *bufconn.Listener
		bufSize   = 1024 * 1024
		ctrl      *gomock.Controller
		mRepo     *mocks.MockIRepository
		mProducer *mocks.MockIProducer
		ctx       context.Context
		conn      *grpc.ClientConn
		client    pb.OcpOfferApiServiceClient
		done      chan struct{}
	)

	BeforeSuite(func() {
		ctrl = gomock.NewController(GinkgoT())
		mRepo = mocks.NewMockIRepository(ctrl)
		mProducer = mocks.NewMockIProducer(ctrl)

		listener = bufconn.Listen(bufSize)
		server := grpc.NewServer()

		pb.RegisterOcpOfferApiServiceServer(server, api.NewOfferAPI(mRepo, mProducer))
		done = make(chan struct{})

		go func() {
			if err := server.Serve(listener); err != nil {
				log.Fatal(err)
			}
			print("done serving")
			<-done
			print("done serving")
			server.Stop()
		}()

	})

	AfterSuite(func() {
		close(done)
		ctrl.Finish()
	})

	BeforeEach(func() {
		ctx = context.Background()
		conn, _ = grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}))
		client = pb.NewOcpOfferApiServiceClient(conn)
	})

	AfterEach(func() {
		conn.Close()
	})

	Context("gRPC call to CreateOfferV1 function", func() {
		When("invalid arguments", func() {
			It("req.UserId, req.Grade, req.TeamId = 0 returns an error codes.InvalidArgument", func() {
				mRepo.EXPECT().
					CreateOffer(gomock.Any(), gomock.Any()).
					Times(0)

				req := &pb.CreateOfferV1Request{UserId: 0, Grade: 0, TeamId: 0}
				res, err := client.CreateOfferV1(ctx, req)

				Expect(res).Should(BeNil())
				Expect(status.Code(err)).Should(BeEquivalentTo(codes.InvalidArgument))
			})
		})

		When("unknown error from CreateOffer", func() {
			It("returns an error", func() {
				mRepo.EXPECT().
					CreateOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(uint64(0), errors.New(""))

				req := &pb.CreateOfferV1Request{UserId: 1, Grade: 2, TeamId: 3}
				res, err := client.CreateOfferV1(ctx, req)

				Expect(res).Should(BeNil())
				Expect(status.Code(err)).Should(BeEquivalentTo(codes.Internal))
			})
		})

		When("normal case", func() {
			It("all props corrected", func() {
				mRepo.EXPECT().
					CreateOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(uint64(1), nil)

				req := &pb.CreateOfferV1Request{UserId: 1, Grade: 2, TeamId: 3}
				res, err := client.CreateOfferV1(ctx, req)

				Expect(res.Id).Should(BeEquivalentTo(1))
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("gRPC call to DescribeOfferV1 function", func() {
		When("invalid arguments", func() {
			It("req.Id = 0 returns an error codes.InvalidArgument", func() {
				mRepo.EXPECT().
					DescribeOffer(gomock.Any(), gomock.Any()).
					Times(0)

				req := &pb.DescribeOfferV1Request{Id: 0}
				res, err := client.DescribeOfferV1(ctx, req)

				Expect(res).Should(BeNil())
				Expect(status.Code(err)).Should(BeEquivalentTo(codes.InvalidArgument))
			})
		})

		When("unknown error from DescribeOffer", func() {
			It("returns an error", func() {
				mRepo.EXPECT().
					DescribeOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New(""))

				req := &pb.DescribeOfferV1Request{Id: 1}
				res, err := client.DescribeOfferV1(ctx, req)

				Expect(res).Should(BeNil())
				Expect(status.Code(err)).Should(BeEquivalentTo(codes.Internal))
			})
		})

		When("normal case", func() {
			It("all props corrected", func() {

				mRepo.EXPECT().
					DescribeOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.Offer{
						ID:     1,
						UserID: 2,
						Grade:  3,
						TeamID: 4,
					}, nil)

				req := &pb.DescribeOfferV1Request{Id: 1}
				res, err := client.DescribeOfferV1(ctx, req)

				Expect(res).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("gRPC call to ListOfferV1 function", func() {
		When("invalid arguments", func() {
			It("initialized values returns an error codes.InvalidArgument", func() {
				mRepo.EXPECT().
					ListOffer(gomock.Any(), gomock.Any()).
					Times(0)

				req := &pb.ListOfferV1Request{
					Pagination: &pb.PaginationInput{},
				}

				res, err := client.ListOfferV1(ctx, req)

				Expect(res).Should(BeNil())
				Expect(status.Code(err)).Should(BeEquivalentTo(codes.InvalidArgument))
			})
		})

		When("unknown error from ListOffer", func() {
			It("returns an error", func() {
				mRepo.EXPECT().
					ListOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, nil, errors.New(""))

				req := &pb.ListOfferV1Request{
					Pagination: &pb.PaginationInput{
						Take: 1,
						Skip: 10,
					},
				}

				res, err := client.ListOfferV1(ctx, req)

				Expect(res).Should(BeNil())
				Expect(status.Code(err)).Should(BeEquivalentTo(codes.Internal))
			})
		})

		When("normal case", func() {
			It("all props corrected", func() {
				offers := make([]models.Offer, 0)

				pagInfo := models.PaginationInfo{
					Page:            1,
					TotalPages:      1,
					TotalItems:      0,
					PerPage:         0,
					HasNextPage:     false,
					HasPreviousPage: false,
				}

				mRepo.EXPECT().
					ListOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(offers, &pagInfo, nil)

				req := &pb.ListOfferV1Request{
					Pagination: &pb.PaginationInput{
						Take: 1,
						Skip: 10,
					},
				}

				res, err := client.ListOfferV1(ctx, req)

				Expect(res).ShouldNot(BeNil())
				Expect(res.Pagination).Should(
					BeEquivalentTo(&pb.PaginationInfo{
						Page:            pagInfo.Page,
						TotalPages:      pagInfo.TotalPages,
						TotalItems:      pagInfo.TotalItems,
						PerPage:         pagInfo.PerPage,
						HasNextPage:     pagInfo.HasNextPage,
						HasPreviousPage: pagInfo.HasPreviousPage,
					}))
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("gRPC call to UpdateOfferV1 function", func() {
		When("invalid arguments", func() {
			It("initialized values returns an error codes.InvalidArgument", func() {
				mRepo.EXPECT().
					UpdateOffer(gomock.Any(), gomock.Any()).
					Times(0)

				req := &pb.UpdateOfferV1Request{
					Id:     0,
					UserId: 0,
					Grade:  0,
					TeamId: 0,
				}

				res, err := client.UpdateOfferV1(ctx, req)

				Expect(res).Should(BeNil())
				Expect(status.Code(err)).Should(BeEquivalentTo(codes.InvalidArgument))
			})
		})

		When("unknown error from RemoveOffer", func() {
			It("returns an error", func() {
				mRepo.EXPECT().
					UpdateOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New(""))

				req := &pb.UpdateOfferV1Request{
					Id:     1,
					UserId: 2,
					Grade:  3,
					TeamId: 4,
				}

				res, err := client.UpdateOfferV1(ctx, req)

				Expect(res).Should(BeNil())
				Expect(status.Code(err)).Should(BeEquivalentTo(codes.Internal))
			})
		})

		When("normal case", func() {
			It("all props corrected", func() {

				mRepo.EXPECT().
					UpdateOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				req := &pb.UpdateOfferV1Request{
					Id:     1,
					UserId: 2,
					Grade:  3,
					TeamId: 4,
				}

				res, err := client.UpdateOfferV1(ctx, req)

				Expect(res).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("gRPC call to RemoveOfferV1 function", func() {
		When("invalid arguments", func() {
			It("initialized values returns an error codes.InvalidArgument", func() {
				mRepo.EXPECT().
					RemoveOffer(gomock.Any(), gomock.Any()).
					Times(0)

				req := &pb.RemoveOfferV1Request{
					Id: 0,
				}

				res, err := client.RemoveOfferV1(ctx, req)

				Expect(res).Should(BeNil())
				Expect(status.Code(err)).Should(BeEquivalentTo(codes.InvalidArgument))
			})
		})

		When("unknown error from RemoveOffer", func() {
			It("returns an error", func() {
				mRepo.EXPECT().
					RemoveOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New(""))

				req := &pb.RemoveOfferV1Request{Id: 1}
				res, err := client.RemoveOfferV1(ctx, req)

				Expect(res).Should(BeNil())
				Expect(status.Code(err)).Should(BeEquivalentTo(codes.Internal))
			})
		})

		When("normal case", func() {
			It("all props corrected", func() {

				mRepo.EXPECT().
					RemoveOffer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				req := &pb.RemoveOfferV1Request{Id: 1}
				res, err := client.RemoveOfferV1(ctx, req)

				Expect(&res).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
	})

})
