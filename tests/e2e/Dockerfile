FROM python:2.7-alpine
RUN apk add --update gcc libffi-dev musl-dev openssl-dev python-dev
WORKDIR /tests/e2e
COPY requirements.txt /tests/e2e
RUN pip install pipenv
RUN pipenv install --system -r requirements.txt --skip-lock --deploy
COPY . /tests/e2e
CMD pytest
