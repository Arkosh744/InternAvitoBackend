FROM golang:1.18

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# build go app
RUN go mod download
RUN go build -o wallet-app ./cmd/server/main.go

CMD ["./wallet-app"]