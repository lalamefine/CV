
FROM golang:latest
WORKDIR /app
COPY . /app
RUN go build -o serve server.go
EXPOSE 80
CMD ["./serve", "mem", "80"]