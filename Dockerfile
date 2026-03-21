# Etapa 1: Compilación
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Descargamos dependencias primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

# Compilamos el binario apuntando a la ruta correcta de tu main.go
# Usamos -ldflags="-s -w" para reducir el tamaño del binario final
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/api/main.go

# Etapa 2: Producción (Imagen final ligera)
FROM alpine:latest

# Certificados necesarios para peticiones HTTPS externas
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Traemos el binario compilado de la etapa anterior
COPY --from=builder /app/main .

# Exponemos el puerto que usa tu app por defecto
EXPOSE 8080

# Ejecutamos la app
CMD ["./main"]
