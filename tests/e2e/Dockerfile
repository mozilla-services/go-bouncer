FROM python:2.7-alpine
RUN apk add --update gcc libffi-dev musl-dev openssl-dev python-dev
WORKDIR /src
RUN pip install pipenv
COPY Pipfile Pipfile.lock /src/
RUN pipenv install --sytem --deploy
COPY . /src
CMD pytest
