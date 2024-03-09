FROM golang:latest

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o /slot-booker

EXPOSE 8081

CMD [ "/slot-booker" ]