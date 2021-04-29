FROM golang:1.16.3-stretch
# create a working directory
WORKDIR /go/app/
# add source code
ADD stup1d-b0t /go/app/

# run main.go
CMD ["/go/app/stup1d-b0t"]