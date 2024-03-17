FROM golang:1.20 as builder

WORKDIR /app


COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o authService ./cmd/authService

FROM alpine:latest
WORKDIR /root/

RUN apk add --no-cache bash && bash --version

COPY --from=builder /app/authService .
COPY --from=builder /app/wait-for-it.sh .
COPY --from=builder /app/openapi .

#COPY --from=builder /app/cert.pem /root/cert.pem
#COPY --from=builder /app/key.pem /root/key.pem

# Сделай скрипт исполняемым
RUN chmod +x wait-for-it.sh

EXPOSE 8080

#CMD ["./authService"]
# Используем wait-for-it.sh для ожидания доступности базы данных перед запуском приложения
CMD ["./wait-for-it.sh", "db:5432", "--", "./authService"]