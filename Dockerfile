FROM golang:1.17-alpine


RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apk update && apk add --no-cache build-base
RUN apk add postgresql-client
RUN chmod +x wait-for-postgres.sh
RUN go mod download && go build -o restapi ./cmd/main.go
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.19.1

CMD ["./restapi"]


