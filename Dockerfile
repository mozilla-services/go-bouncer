FROM golang:1.5

ENV PROJECT github.com/mozilla-services/go-bouncer

COPY . /go/src/$PROJECT

ENV GOPATH /go/src/$PROJECT/Godeps/_workspace:$GOPATH

RUN go install $PROJECT && go install $PROJECT/go-sentry
