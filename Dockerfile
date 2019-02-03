FROM golang:1.11

# Copy app
RUN mkdir /app 
ADD . /app/
WORKDIR /app 

# Build
RUN go build -o main .

# Run
EXPOSE 3000
CMD ["./main"]
