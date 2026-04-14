FROM python:3.10-slim

WORKDIR /app

# Install git and ffmpeg (git is needed for pip install from GitHub)
RUN apt-get update && apt-get install -y \
    git \
    ffmpeg \
    && rm -rf /var/lib/apt/lists/*

RUN pip install --upgrade pip setuptools wheel

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 50053

CMD ["python", "app.py"]