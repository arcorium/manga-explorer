FROM golang:1.21.6 AS builder
LABEL authors="arcorium"

WORKDIR /app

ENV GIN_MODE=release
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . .

RUN go mod download

RUN make build

#FROM builder AS test-runner

#RUN make prepare.test

#RUN go test -v

FROM scratch AS runner

COPY --from=builder /server /server

ENV GIN_MODE=release

EXPOSE 9999

ENTRYPOINT ["/build/server"]