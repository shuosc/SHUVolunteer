#FROM golang:1.12-alpine as builder
#FROM golang:1.12-stretch as builder
#RUN apk add git tesseract-ocr gcc libc-dev g++ libtool libtesseract-dev
#RUN apt install git tesseract-ocr gcc libc-dev g++ libtool libtesseract-dev
FROM jitesoft/tesseract-ocr
RUN apt update
RUN apt -y install golang gcc g++
ENV GO111MODULE=on
COPY . /go/src/SHUVolunteer
WORKDIR /go/src/SHUVolunteer
RUN ls && pwd
RUN go install SHUVolunteer
RUN go get && go build


FROM alpine:latest
COPY --from=builder /go/src/SHUVolunteer/SHUVolunteer /SHUVolunteer
WORKDIR /
CMD ./SHUVolunteer
EXPOSE 8001