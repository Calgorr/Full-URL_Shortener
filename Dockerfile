FROM golang:1.20.4-alpine
WORKDIR /home/app
COPY . ./
RUN go mod tidy
RUN go build -o url_shortener .
CMD ["./url_shortener"]