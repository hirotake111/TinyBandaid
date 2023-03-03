# syntax=docker/dockerfile:1
# Above is the syntax to tell Docker we will use the latest syntax version release.

FROM golang:1.20rc2-alpine3.17 as build

ENV GO111MODULE=on

WORKDIR /app

# Download go modules and dependencies
COPY go.mod ./
RUN go mod download

# Copy all .go files
COPY . .

# RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /serve -a -ldflags '-extldflags "-static"' cmd/lb/main.go

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /serve /serve
EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT [ "/serve" ]
