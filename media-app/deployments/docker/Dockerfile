FROM golang:1.21.6

RUN mkdir -p /build
RUN mkdir -p /app

WORKDIR /build
COPY ./shared ./shared
COPY ./media-app ./media-app

# Download the shared dependencies
WORKDIR /build/shared
RUN go mod download

# Download the app dependencies
WORKDIR /build/media-app
RUN go mod download

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /app/main .

# Remove unnecessary files
RUN rm -rf /build

WORKDIR /app
COPY ./media-app/configs/config.prod.yaml config.yaml

EXPOSE 8004

ENTRYPOINT ["/app/main"]
CMD [ "serve" ]