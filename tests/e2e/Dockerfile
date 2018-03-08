FROM python:2.7-alpine
RUN apk add --update gcc libffi-dev musl-dev openssl-dev python-dev
WORKDIR /src
COPY requirements.txt /src
RUN pip install pipenv==10.1.2
RUN pipenv install --system -r requirements.txt --skip-lock
COPY . /src
CMD pytest
