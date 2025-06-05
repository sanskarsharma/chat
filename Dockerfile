FROM golang:1.14-alpine as builder
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .

ENV GOOS=linux \
    GOARCH=amd64 \
    GOBIN=$GOPATH/bin
RUN go mod download
RUN go build -ldflags="-s -w" -o ./bin/main-bin ./*.go

FROM alpine:3.9

# copying binary built from previous stage
WORKDIR /usr/bin
COPY --from=builder /go/src/app/bin /go/bin

# copy the HTML file to the working directory
COPY --from=builder /go/src/app/chat.html /usr/bin/chat.html

EXPOSE 8080
ENTRYPOINT /go/bin/main-bin 