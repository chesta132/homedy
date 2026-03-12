# ================================
# Stage 1: Build OS
# ================================
FROM ubuntu:18.04 AS buildos

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    wget \
    git \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# ================================
# Stage 2: Build
# ================================
FROM golang:1.25.5 AS build

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

# Copy source code
COPY . .

RUN go mod download

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -o main .

# ================================
# Stage 3: Run
# ================================
FROM ubuntu:18.04 AS run

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]