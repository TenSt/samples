FROM golang:1.13.4-stretch

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 

RUN go get -u "github.com/gorilla/mux"

RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -o main main.go

EXPOSE 8080

CMD ["/app/main"]