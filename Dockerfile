FROM alpine as certs

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

FROM golang:1.17 as builder
WORKDIR /go/src/node-pool-shifter
ADD . /go/src/node-pool-shifter

RUN go get -d -v ./...

RUN go build -o /go/bin/estafette-gke-node-pool-shifter-v2

FROM gcr.io/distroless/base

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/estafette-gke-node-pool-shifter-v2 /

CMD ["/estafette-gke-node-pool-shifter-v2"]
