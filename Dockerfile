FROM golang:1.8

ENV PROJECT=github.com/mozilla-services/go-bouncer

COPY version.json /app/version.json
COPY . /go/src/$PROJECT

EXPOSE 8000

RUN go install $PROJECT && go install $PROJECT/go-sentry
CMD ["go-bouncer", "--addr", "127.0.0.1:8000"]
