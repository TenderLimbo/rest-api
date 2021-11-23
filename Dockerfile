FROM golang:1.17-alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apk update && apk add --no-cache git
RUN apk add postgresql-client
RUN chmod +x wait-for-postgres.sh
RUN go mod download && go build -o restapi ./cmd/main.go

CMD ["./restapi"]