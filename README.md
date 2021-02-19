# Hea Meelega - CRM for happy You

## Local development

### Run DB container

     docker run --rm --name postgres -e POSTGRES_DB=wallester -e POSTGRES_USER=wallester -e POSTGRES_PASSWORD=12345 -p 5432:5432 postgres:13.1-alpine

### Build and run app container

    docker build -f .infra/Dockerfile -t heameelega:latest .
    docker run --rm -it --link postgres:postgres -p 80:80 heameelega
