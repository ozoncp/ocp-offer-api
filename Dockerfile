# Builder

FROM golang:1.16-alpine AS builder

RUN apk add --update make git protoc protobuf protobuf-dev

COPY . /home/github.com/ozoncp/ocp-offer-api

WORKDIR /home/github.com/ozoncp/ocp-offer-api

RUN make deps && make build


# gRPC Server

FROM alpine:latest as server

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /home/github.com/ozoncp/ocp-offer-api/bin/grpc-server .
COPY --from=builder /home/github.com/ozoncp/ocp-offer-api/migrations/ .

RUN chown root:root grpc-server

EXPOSE 50051
EXPOSE 8080
EXPOSE 9100

CMD ["./grpc-server"]


# Kafka consumer

FROM alpine:latest as consumer

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /home/github.com/ozoncp/ocp-offer-api/bin/kafka-consumer .

RUN chown root:root kafka-consumer

CMD ["./kafka-consumer"]