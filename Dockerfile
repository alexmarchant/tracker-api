FROM golang:1.11

# Get watcher lib for dev
RUN go get github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

# Copy app
WORKDIR /app
COPY . .

# Build
RUN go build -o main .

# Run
EXPOSE 3000
CMD ["./main"]
