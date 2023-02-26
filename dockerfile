FROM golang:1.18

WORKDIR /app

# Install air
RUN go install github.com/cosmtrek/air@latest
RUN apt install -y tzdata

# Warm go mod caches
ADD go.mod go.mod
ADD go.sum go.sum
RUN go mod download

CMD ["air", "-c", ".air.toml"]