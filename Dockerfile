FROM golang:1.21.6 AS builder
LABEL authors="arcorium"

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . .

RUN go mod download

RUN make build

RUN make migrate

#FROM builder AS test-runner

#RUN make prepare.test

#RUN go test -v

FROM alpine AS runner

COPY --from=builder /app/build/server /app/
COPY --from=builder /app/.env /app/

WORKDIR /app

ENV GIN_MODE=release

EXPOSE 9999

ENTRYPOINT ["./server"]