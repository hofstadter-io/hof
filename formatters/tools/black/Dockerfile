FROM python:3.9-alpine

WORKDIR /app

COPY Pipfile Pipfile.lock /app/

RUN apk add --no-cache --virtual .build-deps build-base git && \
    pip install --upgrade pip pipenv && \
    pipenv install --deploy --system && \
    pip uninstall -y pipenv && \
    apk del .build-deps && \
    rm -rf /root/.cache

COPY . /app

EXPOSE 3000
CMD ["gunicorn", "app:app", "--bind", "0.0.0.0:3000", "--log-file", "-"]
