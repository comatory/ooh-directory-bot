FROM golang:1.23-alpine

WORKDIR /app
COPY . /app
VOLUME data

RUN apk add --no-cache make

RUN go build -o /app/ooh-my-bot .

RUN echo "0 10 * * * /app/ooh-my-bot --config-file=/app/data/.env --records-file=/app/data/records.txt  >> /app/data/log.txt" > /etc/crontabs/root

CMD ["crond", "-f", "-d", "8"]
