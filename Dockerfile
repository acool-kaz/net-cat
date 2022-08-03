FROM golang:latest
WORKDIR /app
COPY . /app/
RUN go build -o main server/server.go
CMD [ "./main" ]