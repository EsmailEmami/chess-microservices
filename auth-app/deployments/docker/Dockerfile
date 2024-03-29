FROM golang:1.21.6

RUN mkdir -p /build
RUN mkdir -p /app
RUN mkdir -p /app/migrations

WORKDIR /build
COPY ./shared ./shared
COPY ./auth-app ./auth-app

# Download the shared dependencies
WORKDIR /build/shared
RUN go mod download

# Download the app dependencies
WORKDIR /build/auth-app
RUN go mod download

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /app/main .

# Remove unnecessary files
RUN rm -rf /build

WORKDIR /app
COPY ./auth-app/configs/config.prod.yaml config.yaml

EXPOSE 8001

ENTRYPOINT ["/app/main"]
CMD [ "serve" ]