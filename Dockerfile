FROM golang:1.6

ENV PROJECT=github.com/mozilla-services/go-bouncer

COPY . /go/src/$PROJECT

RUN go install $PROJECT && go install $PROJECT/go-sentry
