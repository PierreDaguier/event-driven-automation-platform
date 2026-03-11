FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/worker ./cmd/worker

FROM gcr.io/distroless/static-debian12:nonroot AS server
COPY --from=builder /out/server /server
ENTRYPOINT ["/server"]

FROM gcr.io/distroless/static-debian12:nonroot AS worker
COPY --from=builder /out/worker /worker
ENTRYPOINT ["/worker"]
