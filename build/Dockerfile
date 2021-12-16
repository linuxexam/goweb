# context should be the repo root
FROM golang:1.17

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -v -o ./goweb

FROM alpine:latest
WORKDIR /app
COPY --from=0 /go/src/app/goweb ./

EXPOSE 8080
CMD ["cd /app/; ./goweb"]
