FROM golang:1.22 AS builder

WORKDIR /usr/src/builder

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o dist/server

FROM ubuntu:24.04 AS server

WORKDIR /usr/src/app

COPY --from=builder /usr/src/builder/dist .
COPY --from=builder /usr/src/builder/dos-server-sdk.json ./

RUN apt update && apt install ca-certificates -y

ENV GOOGLE_APPLICATION_CREDENTIALS=/usr/src/app/dos-server-sdk.json

CMD [ "./server" ]
