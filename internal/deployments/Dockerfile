# 1단계: go 빌드
FROM golang:1.23.0 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy
COPY . .
# 바이너리 빌드
RUN go build -o automate_youtube_subscription ./cmd/automate_youtube_subscription
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o automate_youtube_subscription ./cmd/automate_youtube_subscription

# 2단계: 실행
FROM debian:bookworm-slim

WORKDIR /app

# 인증 관련 패키지가 필요하다면 설치 (선택적)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# 빌드된 바이너리만 복사
COPY --from=builder /app/automate_youtube_subscription .
COPY .env .env
COPY client_secret*.json client_secret.json

# 실행
ENTRYPOINT ["./automate_youtube_subscription"]