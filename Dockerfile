FROM golang
RUN apt update
RUN apt -y install golang gcc g++ git libleptonica-dev tesseract-ocr libtesseract-dev > /dev/null
COPY . /root/go/src/SHUVolunteer
WORKDIR /root/go/src/SHUVolunteer
RUN go get && go build
CMD ./SHUVolunteer
EXPOSE 8001
