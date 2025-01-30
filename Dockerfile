# Etapa 1: Construir el binario
FROM golang:1.23.5-alpine3.21

# Directorio de trabajo
WORKDIR /app

# Copiar dependencias y codigo fuente
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Compilar el binario
RUN go build -o scraper-service cmd/main.go


# Etapa 2: Ejecutar el binario
FROM alpine:latest

# Establecer el directorio de trabajo
WORKDIR /root/

# Instalar openssl para manejar certificados y certificados CA
RUN apk add --no-cache openssl ca-certificates && update-ca-certificates

# Copiar el binario de la etapa 1
COPY --from=0 /app/scraper-service .

# Copiar el certificado de la base de datos
COPY --from=0 /app/internal/config/shipdb.pem .

# Ejecutar el binario
CMD ["./scraper-service"]