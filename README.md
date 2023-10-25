# Onboarding Go Project

## Installation

This project contains a Dockerfile and docker-compose.yml. To build and launch the application with its dependencies,
run:

``` 
docker-compose up --build
```

## Test

If application is running, stop it with:

```
docker stop ps-tag-onboarding-go-app-1 ps-tag-onboarding-go-mongo-1
```

And then run tests with:

```
cd tests/
docker-compose -f docker-compose-test.yml up -d
go test
```