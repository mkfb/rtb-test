FROM python:3.9-slim-buster

WORKDIR /app

COPY traffic/generator.py /app/generator.py

RUN pip install aiohttp

CMD ["python", "/app/generator.py", "http://nginx/bid"]