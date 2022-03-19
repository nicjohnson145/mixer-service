FROM golang:1-alpine AS builder
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR /src/cmd/mixer-server
RUN CGO_ENABLED=0 go build -buildvcs=false

FROM scratch
COPY --from=builder /src/cmd/mixer-server/mixer-server .
CMD ["./mixer-server"]
