#FROM golang:1.12-alpine as builder
FROM golang:1.12-stretch as builder
#RUN apk add git tesseract-ocr gcc libc-dev g++ libtool libtesseract-dev
RUN apt install git tesseract-ocr gcc libc-dev g++ libtool libtesseract-dev
ENV GO111MODULE=on
COPY . /go/src/shu-volunteer
WORKDIR /go/src/shu-volunteer
RUN ls
RUN go get && go build


FROM alpine:latest
COPY --from=builder /go/src/shu-volunteer/shu-volunteer /shu-volunteer
WORKDIR /
CMD ./shu-volunteer
EXPOSE 8001