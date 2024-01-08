FROM golang:1.20-bullseye

ENV PROJECT=github.com/mozilla-services/go-bouncer

COPY . /app

EXPOSE 8000

WORKDIR /app

RUN go install -mod vendor $PROJECT

CMD ["go-bouncer", "--addr", "127.0.0.1:8000"]
