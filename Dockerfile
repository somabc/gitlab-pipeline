FROM golang:alpine as builder

WORKDIR /go/src/bcochran/gitlab-pipeline

COPY . .

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

RUN CGO_ENABLED=0 GOOS=linux go build -a  -installsuffix cgo -o main .


FROM alpine

WORKDIR /trigger

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/src/bcochran/gitlab-pipeline/main .

CMD ["./main"]
