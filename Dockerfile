FROM golang:1.9-alpine as builder

RUN apk add --update --no-cache ca-certificates

# build
WORKDIR /go/src/github.com/monostream/muescheli/

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o muescheli .

# copy artefacts
WORKDIR /app/
RUN cp /go/src/github.com/monostream/muescheli/muescheli .
RUN rm -r /go/src/

FROM scratch

WORKDIR /app/
COPY --from=builder /app/ .

# copy certificates so that files can be fetched from ssl sites
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# setup environment
ENV PATH "/app:${PATH}"

CMD ["muescheli"]
