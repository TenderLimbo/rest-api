FROM golang:1.17-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client
RUN chmod +x wait-for-postgres.sh
RUN go mod download && go build -o restapi ./cmd/main.go

CMD ["./restapi"]