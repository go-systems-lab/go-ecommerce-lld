FROM python:3.12-slim

WORKDIR /app

RUN apt-get update && apt-get install -y build-essential libpq-dev && rm -rf /var/lib/apt/lists/*

COPY recommender/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY recommender /app

ENV PYTHONPATH="/app:${PYTHONPATH}"
