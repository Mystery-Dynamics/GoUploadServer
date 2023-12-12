FROM golang:latest AS build
COPY . /APP
WORKDIR /APP
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o upload-server

FROM alpine:latest AS product
COPY --from=build /APP/upload-server /APP/upload-server
COPY --from=build /APP/templates /APP/templates
WORKDIR /APP
ENTRYPOINT ["/APP/upload-server"]
