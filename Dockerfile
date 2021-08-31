FROM golang:alpine as builder
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY estafette-gke-node-pool-shifter-v2 /

CMD ["./estafette-gke-node-pool-shifter-v2"]
