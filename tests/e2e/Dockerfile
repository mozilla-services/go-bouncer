FROM python:2.7-alpine
RUN apk --no-cache add gcc
WORKDIR /src
COPY requirements.txt /src
RUN pip install pipenv
RUN pipenv install --system -r requirements.txt --skip-lock --deploy
COPY . /src
CMD pytest
