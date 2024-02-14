FROM golang:1.21.2-alpine as builder
RUN apk update
RUN apk add --no-cache gcc musl-dev

WORKDIR /usr/src/app
COPY . .

ENV GO111MODULE=on

# RUN go get -u github.com/off-chain-storage/go-off-chain-storage@v1.2.1
# RUN go get -u github.com/ethereum/go-ethereum@latest
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o bin/curie-node cmd/curie-node/main.go
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o bin/proposer cmd/proposer/main.go


### Executable Image
FROM alpine

RUN apk add --no-cache libstdc++

COPY --from=builder /usr/src/app/bin/curie-node ./curie-node
COPY --from=builder /usr/src/app/bin/proposer ./proposer


# ENTRYPOINT ["./curie-node"]