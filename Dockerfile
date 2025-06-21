FROM golang:1.24.4-alpine AS builder

WORKDIR $GOPATH/src/github.com/satisfactorymodding/smr-api/

ENV GO111MODULE=on

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o /go/bin/api cmd/api/serve.go


FROM golang:1.24.4-alpine
COPY --from=builder /go/bin/api /api
WORKDIR /app
COPY static /app/static
COPY migrations /app/migrations
EXPOSE 5020
ENTRYPOINT ["/api"]
