# Hea Meelega - CRM for happy You

![](https://github.com/dmitriivoitovich/heameelega/workflows/build/badge.svg)

## Local development

### Run DB container

     docker run --rm --name postgres -e POSTGRES_DB=heameelega -e POSTGRES_USER=heameelega -e POSTGRES_PASSWORD=12345 -p 5432:5432 postgres:13.1-alpine

### Build and run app container

    docker build -f .infra/docker/Dockerfile -t dmitriivoitovich/heameelega:latest .
    docker run --rm -it --link postgres:postgres -e VIPER_DB_HOST=postgres -p 80:80 heameelega

### Run static code analyser

    brew install golangci-lint
    golangci-lint run
