FROM golang:1.8.3

RUN go get github.com/Masterminds/glide

RUN mkdir -p /go/src/github.com/jobtalk/hawkeye

WORKDIR /go/src/github.com/jobtalk/hawkeye

COPY . /go/src/github.com/jobtalk/hawkeye

RUN glide i

RUN go build . \
	&& cp hawkeye /bin/hawkeye

CMD hawkeye
