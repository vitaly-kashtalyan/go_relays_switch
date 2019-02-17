FROM golang:1.11.5

# install packeages
RUN go get -d -v github.com/gin-contrib/cors
RUN go get -d -v github.com/gin-contrib/gzip
RUN go get -d -v github.com/gin-contrib/location
RUN go get -d -v github.com/gin-gonic/gin
RUN go get -d -v github.com/vitaly-kashtalyan/hlk-sw16

# create a working directory
WORKDIR /go/src/app
# add source code
COPY main.go main.go
# build main.go
RUN go build ./main.go
EXPOSE 8082
# run the binary
CMD ["./main"]