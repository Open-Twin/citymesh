FROM python:3

RUN pip install kafka-python

ADD app /app/

COPY . /app/

CMD ["python3", "-u", "app/producer.py"]
