FROM alpine:3.12

# run: GOOS=linux GOARCH=amd64 go build -o easyserver .
COPY easyserver /usr/bin/easyserver

CMD ["/usr/bin/easyserver", "serve", "/data"]
