FROM golang:1.4

RUN go get github.com/tools/godep
COPY . /go/src/github.com/cloudnautique/vol-cleanup
WORKDIR /go/src/github.com/cloudnautique/vol-cleanup
RUN godep go build -a -tags "netgo" -installsuffix netgo -ldflags "-extldflags -static" -o /usr/bin/vol-cleanup
CMD ["vol-cleanup", "-h"]
