FROM golang:1.17.0-alpine as builder
WORKDIR /go/src/app
ADD . .

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

RUN go get -d -v ./...

RUN go build -o /estafette-gke-node-pool-shifter-v2

FROM gcr.io/distroless/base

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /estafette-gke-node-pool-shifter-v2 /

CMD ["./estafette-gke-node-pool-shifter-v2"]
