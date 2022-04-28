# Stage 1
FROM golang:1.18

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN mkdir -p /usr/local/bin && CGO_ENABLED=0 go build -v -o /usr/local/bin/app .

# Stage 2
FROM scratch

COPY --from=0 /usr/local/bin/app ./

CMD ["./app"]
