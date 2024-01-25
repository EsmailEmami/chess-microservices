FROM golang:1.20

RUN mkdir -p /build
RUN mkdir -p /app

WORKDIR /build
COPY . .

# Download the shared dependencies
WORKDIR /build/shared
RUN go mod download

# Download the app dependencies
WORKDIR /build/game-app
RUN go mod download

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /app/main ./cmd

# Remove unnecessary files
RUN rm -rf /build

WORKDIR /app
COPY ./game-app/configs/config.prod.yaml config.yaml

EXPOSE 8002

ENTRYPOINT ["/app/main"]
CMD [ "serve" ]