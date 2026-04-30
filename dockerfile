FROM golang:1.26-alpine

# Instala ferramentas de compilação e dependências de sistema
RUN apk add --no-cache \
    build-base \
    libaio \
    libnsl \
    libc6-compat \
    unzip \
    wget

# Baixa e instala o Oracle Instant Client (versão básica)
WORKDIR /opt/oracle
RUN wget https://download.oracle.com/otn_software/linux/instantclient/2114000/instantclient-basiclite-linux.x64-21.14.0.0.0dbru.zip && \
    unzip instantclient-basiclite-linux.x64-21.14.0.0.0dbru.zip && \
    rm instantclient-basiclite-linux.x64-21.14.0.0.0dbru.zip && \
    mv instantclient_21_14 instantclient

# Configura as variáveis de ambiente para o Oracle encontrar as bibliotecas
ENV LD_LIBRARY_PATH=/opt/oracle/instantclient
ENV CGO_ENABLED=1

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

EXPOSE 8080

CMD ["./main"]