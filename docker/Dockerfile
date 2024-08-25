# stage 1: build stellar-etl app
FROM golang:1.23.0-alpine AS build

WORKDIR /usr/src/etl

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin ./...

# stage 2: runtime enviroment
FROM stellar/stellar-core:21.3.1-2007.4ede19620.focal

WORKDIR /etl

COPY --from=build /usr/local/bin/stellar-etl /usr/local/bin/stellar-etl
COPY --from=build /usr/src/etl/docker docker

# changing workdir to a new path in order to use mounted empty ephemeral volumes as storage
WORKDIR /etl/data

# clear entrypoint from stellar-core image
ENTRYPOINT []

CMD ["stellar-etl"]
