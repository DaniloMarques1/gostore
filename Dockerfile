FROM golang:1.16

WORKDIR /gostore

COPY . .

RUN go mod download
RUN go mod tidy

RUN go build .

EXPOSE 5000

CMD ["./gostore"]
