## How to run with Docker

1. Build the Docker image
```bash
docker build -t itu-minitwit-cs .
```

2. Run the Docker container
```bash
docker run -d -p 8080:8080 -v <path_to_minitwit_db>:/app/data/minitwit.db itu-minitwit-cs
```