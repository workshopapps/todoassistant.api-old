FROM golang:1.19-alpine3.16

RUN mkdir -p /home/backend

COPY . /home/backend/

WORKDIR /home/backend/

RUN go build -o main

EXPOSE  2022

ENTRYPOINT [ "go", "run", "main.go" ]