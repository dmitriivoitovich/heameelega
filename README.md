# Hea Meelega - CRM for happy You

## Local development

### Run DB container

     docker run --rm --name postgres -e POSTGRES_DB=heameelega -e POSTGRES_USER=heameelega -e POSTGRES_PASSWORD=12345 -p 5432:5432 postgres:13.1-alpine

### Build and run app container

    docker build -f .infra/Dockerfile -t heameelega:latest .
    docker run --rm -it --link postgres:postgres -e VIPER_DB_HOST=postgres -p 80:80 heameelega

### Run static code analyser

    brew install golangci-lint
    golangci-lint run -E asciicheck,bodyclose,cyclop,depguard,dogsled,dupl,durationcheck,errorlint,exhaustive,exportloopref,forbidigo,funlen,gocognit,gocritic,gocyclo,godot,godox,goerr113,gofmt,goheader,goimports,golint,gomnd,gomodguard,goprintffuncname,gosec,ifshort,interfacer,makezero,maligned,misspell,nakedret,nestif,nlreturn,noctx,nolintlint,paralleltest,prealloc,predeclared,revive,rowserrcheck,scopelint,sqlclosecheck,stylecheck,testpackage,thelper,tparallel,unconvert,unparam,whitespace,wsl
