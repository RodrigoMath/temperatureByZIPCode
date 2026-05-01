# ========================================
# Stage 1: Build
# ========================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# ========================================
# Stage 2: Run
# ========================================
FROM alpine:latest

# Instalar certificados CA para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binário compilado
COPY --from=builder /app/main .

# Expor porta
EXPOSE 8080

# Variável de ambiente padrão para porta
ENV PORT=8080

# Comando de execução
CMD ["./main"]
