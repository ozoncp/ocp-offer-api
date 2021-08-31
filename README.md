# Ozon Code Platform Offer API

Information about the issued offer to the student

by [Evald Smalyakov](https://github.com/evald24)

---

## Build project

### Local

For local assembly you need to perform

```zsh
$ make deps # Installation of dependencies
$ make build # Build project
```

### Based on tag

When publishing a tag, there is an automatic build and the publication of it on Github Packages where the tag number is version docker image.

- `builder` - The image in which the project build happens is also used for cache and saving time when building

- `kafka-consumer` - The docker image, consumer which asynchronous writes in the database

- `grpc-server` - The docker image, сontains basic business logic

---

## Running

### For local development

```zsh
$ docker-compose up -d
```

### Other

For example, launching an image based on a release

```zsh
$ docker-compose -f docker-compose.stage.yml up -d
```

---

## Services

### Swagger UI

The Swagger UI is an open source project to visually render documentation for an API defined with the OpenAPI (Swagger) Specification

- http://localhost:9080

### Grafana:

- http://localhost:3000/
- - login `admin`
- - password `admin`

### gRPC:

- http://localhost:50051/

### Gateway:

It reads protobuf service definitions and generates a reverse-proxy server which translates a RESTful HTTP API into gRPC

- http://localhost:8080

### Metrics:

Metrics GRPC Server

- localhost:9100/metrics

### Status:

Service condition and its information

- http://ocp-offer-api.evaldsmalyakov.dev:8000/
- - `/live`- Layed whether the server is running
- - `/ready` - Is it ready to accept requests
- - `/version` - Version and assembly information

### Prometheus:

Prometheus is an open-source systems monitoring and alerting toolkit

- http://localhost:9090/

### Kafka

Apache Kafka is an open-source distributed event streaming platform used by thousands of companies for high-performance data pipelines, streaming analytics, data integration, and mission-critical applications.

- http://localhost:9094
- http://kafka:9092/

### Kafka UI

UI for Apache Kafka is a simple tool that makes your data flows observable, helps find and troubleshoot issues faster and deliver optimal performance. Its lightweight dashboard makes it easy to track key metrics of your Kafka clusters - Brokers, Topics, Partitions, Production, and Consumption.

- http://localhost:9001

### Jaeger UI

Monitor and troubleshoot transactions in complex distributed systems.

- http://localhost:16686

### Graylog

Graylog is a leading centralized log management solution for capturing, storing, and enabling real-time analysis of terabytes of machine data.

- http://localhost:9000
- - login `admin`
- - password `admin`
