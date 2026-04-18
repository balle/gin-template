FROM golang:1.25-alpine

RUN apk update
RUN apk upgrade
RUN apk add git

RUN mkdir /code
COPY . /code/
WORKDIR /code

ENTRYPOINT [ "go", "run", "main.go" ] 

EXPOSE 8000/tcp
