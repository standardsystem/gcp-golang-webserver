FROM golang:latest as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/gcp-cloudrun-gcs
COPY . .
RUN go build -o webserver

# runtime image
FROM alpine
COPY --from=builder /go/src/gcp-cloudrun-gcs /app

CMD /app/webserver $PORT