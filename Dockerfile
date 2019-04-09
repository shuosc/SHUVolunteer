#FROM golang:1.12-alpine as builder
#FROM golang:1.12-stretch as builder
#RUN apk add git tesseract-ocr gcc libc-dev g++ libtool libtesseract-dev
#RUN apt install git tesseract-ocr gcc libc-dev g++ libtool libtesseract-dev
FROM jitesoft/tesseract-ocr as builder
RUN apt update
RUN apt -y install golang gcc g++ git libtesseract-dev
COPY . /root/go/src/SHUVolunteer
WORKDIR /root/go/src/SHUVolunteer
RUN ls && pwd && echo $GOPATH && echo $GOROOT
RUN go get && go build


FROM alpine:latest
COPY --from=builder /root/go/src/SHUVolunteer/SHUVolunteer /SHUVolunteer
WORKDIR /
CMD ./SHUVolunteer
EXPOSE 8001