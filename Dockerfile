FROM golang:1.15.2-alpine3.12
WORKDIR /go/src/app
COPY . .
RUN go install -mod=vendor cmd/server/lgtm.go
CMD lgtm