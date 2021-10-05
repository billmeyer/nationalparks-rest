# https://docs.docker.com/language/golang/build-images/

##
## BUILD
##

FROM golang:alpine AS build
ARG GOARCH="amd64"

RUN apk update && \
    apk add curl \
            git \
            bash \
    rm -rf /var/cache/apk/*

WORKDIR /go/src/nationalparks-rest

# copy module files first so that they don't need to be downloaded again if no change
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

# copy source files and build the binary
COPY . ./
RUN go build cmd/main/server.go

##
## DEPLOY
##

FROM alpine:latest
WORKDIR /
COPY --from=build /go/src/nationalparks-rest/server /nationalparks-rest
EXPOSE 8080
ENTRYPOINT ["/nationalparks-rest"]
