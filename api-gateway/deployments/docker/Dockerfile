FROM golang:1.21.6

RUN mkdir -p /build
RUN mkdir -p /app

WORKDIR /build
COPY ./shared ./shared
COPY ./api-gateway ./api-gateway

# Download the shared dependencies
WORKDIR /build/shared
RUN go mod download

# Download the app dependencies
WORKDIR /build/api-gateway
RUN go mod download

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /app/main ./cmd

# Remove unnecessary files
RUN rm -rf /build

WORKDIR /app
COPY ./api-gateway/configs/app/config.prod.yaml config.yaml
COPY ./api-gateway/configs/gateway/gateway.production.json .

EXPOSE 8000

ENTRYPOINT ["/app/main"]