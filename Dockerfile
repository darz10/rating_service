FROM golang:1.17-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make run-server.sh executable
RUN chmod +x run-server.sh

# build go app
RUN go mod download
RUN go build -o ratingBookService ./cmd/main.go

CMD ["./ratingBookService"]