FROM golang:1.25.4-alpine AS builder 
WORKDIR /app

COPY go.mod go.sum ./ 
RUN --mount=type=cache,target=/go/pkg/mod go mod download 

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/manager/

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/main .
COPY migrations ./migrations
COPY *.env .
EXPOSE 8080
ENTRYPOINT [ "./main" ]

