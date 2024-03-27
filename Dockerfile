FROM golang:1.22

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y \
    && rm -rf /var/lib/apt/lists/*

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o jeopardy-data-scraper ./cmd/jeopardy-data-scraper

# Defaults
ENV APP_MODE=""
ENV DB_HOST=""
ENV DB_USERNAME=""
ENV DB_PASSWORD=""
ENV DB_NAME=""
ENV DB_PORT="5432"
ENV DB_TIMEZONE="America/New_York"

CMD ["./jeopardy-data-scraper"]