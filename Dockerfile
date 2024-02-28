FROM golang:1.22

COPY . /forum

WORKDIR /forum

RUN apt update && apt install -y sqlite3
RUN go mod download
RUN go run -o forum cmd/forum/main.go

CMD ["./forum/forum"]