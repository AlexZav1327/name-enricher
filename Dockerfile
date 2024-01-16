FROM golang:latest as builder
ADD . /src/app
WORKDIR /src/app
RUN CGO_ENABLED=0 GOOS=linux go build -o enricher-service ./cmd/enricher-service/main.go
EXPOSE 8086

FROM alpine:edge
COPY --from=builder /src/app/enricher-service /enricher-service
RUN chmod +x ./enricher-service
ENTRYPOINT ["/enricher-service"]