# builder1
FROM golang:1.16 as builder1
## init setting
WORKDIR /workdir
ENV GO111MODULE="on"
ARG VERSION
## download packages
COPY go.mod go.sum ./
RUN go mod download
## build
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X \"main.appVersion=${VERSION}\"" -o app

# runner
FROM alpine:latest as runner
## copy binary
COPY --from=builder1 /workdir/app .
## Run
EXPOSE 8080
ENTRYPOINT ["./app"]
