FROM python:2.7-alpine
WORKDIR /tests/e2e
COPY requirements.txt /tests/e2e
RUN pip install pipenv
RUN pipenv install --system -r requirements.txt --skip-lock --deploy
COPY . /tests/e2e
CMD pytest
