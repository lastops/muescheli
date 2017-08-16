FROM golang:1.8-alpine as builder

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

# setup environment
ENV PATH "/app:${PATH}"

CMD ["./muescheli"]
