FROM golang:1.21.6 AS builder
LABEL authors="arcorium"

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . .

RUN go mod download

RUN make build

#FROM builder AS test-runner

#RUN make prepare.test

#RUN go test -v

FROM alpine:3.19 AS runner

COPY --from=builder /app/build/ /app/.env /app/

WORKDIR /app

ENV GIN_MODE=release

EXPOSE 9999

ENTRYPOINT ["sh", "-c", "./migrate --env --special --no-ssl && ./server"]