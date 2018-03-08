FROM python:2.7-alpine
RUN apk add --update gcc libffi-dev musl-dev openssl-dev python-dev
WORKDIR /src
COPY requirements.txt pipenv.txt /src/
RUN pip install -r pipenv.txt
RUN pipenv install --system -r requirements.txt --skip-lock
COPY . /src
CMD pytest
