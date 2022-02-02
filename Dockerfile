FROM golang:1.17-alpine

WORKDIR /app

COPY app/go.mod ./
COPY app/go.sum ./
RUN go mod download

COPY app/ ./

RUN go build -o router

EXPOSE 8080

CMD [ "./router" ]