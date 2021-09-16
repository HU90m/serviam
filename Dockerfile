FROM golang:1.17-alpine

WORKDIR /go/src/serviam

COPY . .
#COPY ./serviam.go .
#COPY ./go.mod .
#COPY ./common .
#COPY ./structs .
#COPY ./files .
#COPY ./internal .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["serviam"]
