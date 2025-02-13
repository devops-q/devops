> [!NOTE]  
> When running locally without Docker, remember to set ASPNETCORE_ENVIRONMENT=Development so the
> appsettings.Development.json config is loaded and the database is read
> from "../minitwit.db" instead of "/app/data/minitwit.db".

## How to run with Docker

1. Build the Docker image

```bash
docker build -t itu-minitwit-cs .
```

2. Run the Docker container

```bash
docker run -d -p 8080:8080 -v <path_to_minitwit_db>:/app/data/minitwit.db itu-minitwit-cs
```