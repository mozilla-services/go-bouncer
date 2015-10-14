FROM golang:1.5

ENV PROJECT=github.com/mozilla-services/go-bouncer \
    GOPATH=/go/src/github.com/mozilla-services/go-bouncer/Godeps/_workspace:$GOPATH

COPY . /go/src/$PROJECT

RUN go install $PROJECT && go install $PROJECT/go-sentry
