FROM golang:1.18-alpine as builder
WORKDIR /build
ADD . /build
RUN go build

FROM alpine
COPY --from=builder /build/ceobe-bot /usr/local/bin
CMD ["ceobe-bot"]
