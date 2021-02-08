FROM golang:1.13.0-alpine

RUN mkdir /go/src/MovieVote
ADD . /go/src/MovieVote
WORKDIR /go/src/MovieVote

RUN go install

ENV PORT=8080

RUN cd $GOPATH/bin && pwd

CMD ["/go/bin/MovieVote"]