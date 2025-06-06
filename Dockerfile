# base image
FROM golang:1.24

# set working dir
WORKDIR /app

# copy and download go modules
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# build binary from main package in target dir
RUN go build -o server ./cmd/server

# expose port to get mapped later
EXPOSE 8080

CMD ["./server"]