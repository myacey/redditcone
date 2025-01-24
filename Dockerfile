FROM golang:1.23-alpine
RUN apk add --no-cache netcat-openbsd
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o redditclone ./cmd/redditclone/main.go

RUN chmod +x ./entrypoint.sh

EXPOSE 8080

ENTRYPOINT [ "./entrypoint.sh" ]
CMD [ "/app/redditclone" ]
