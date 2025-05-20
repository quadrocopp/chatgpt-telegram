FROM golang:1-alpine as builder
RUN apk update && apk add make
WORKDIR /build
ADD . .
RUN make build

FROM alpine
ENV PORT=8080
EXPOSE 8080
COPY --from=builder /build/chatgpt-telegram /bin/chatgpt-telegram
RUN chmod +x /bin/chatgpt-telegram

ENTRYPOINT ["/bin/chatgpt-telegram"]