FROM golang:alpine

RUN go install github.com/ddmendes/hello@latest
WORKDIR /go/src/app
EXPOSE 8080

CMD ["hello"]

