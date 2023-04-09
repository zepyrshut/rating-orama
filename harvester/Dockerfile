FROM python:3.11.3-alpine
LABEL authors="Pedro Pérez"

COPY requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["python", "main.py"]