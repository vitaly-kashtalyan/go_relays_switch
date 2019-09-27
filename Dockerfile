FROM golang:1.13-alpine
WORKDIR /app
ARG HLK_SW16_HOST=192.168.0.200
ARG HLK_SW16_PORT=8080
ARG GIN_MODE=release
ARG APP_PORT=8082
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8082
CMD ["./main"]