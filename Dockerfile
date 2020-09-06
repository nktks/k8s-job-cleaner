FROM golang:1.15-alpine3.12 as builder

WORKDIR /go/src/github.com/nakatamixi/k8s-job-cleaner

COPY go.mod go.sum  ./
RUN apk add --no-cache git bash &&\
    go mod download

COPY . .
RUN go build -v -i -o bin/app ./

FROM alpine:3.12

COPY --from=builder /go/src/github.com/nakatamixi/k8s-job-cleaner/bin/app .
RUN apk add --no-cache ca-certificates
ENTRYPOINT ["./app"]
