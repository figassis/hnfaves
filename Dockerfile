FROM golang:alpine as builder
ENV GO111MODULE on
ENV APP_PORT 8080
ENV REQUESTS_PER_SECOND 3

RUN mkdir -p /build

WORKDIR /build
COPY go.mod go.sum ./
RUN apk add git && go mod download
ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /app .

# FROM alpine
FROM gcr.io/distroless/base
COPY --from=builder /app /app

WORKDIR /

RUN mkdir -p /cache
VOLUME [ "/cache" ]

EXPOSE 80 8080

CMD ["./app"]