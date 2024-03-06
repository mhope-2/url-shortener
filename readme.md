# URL Shortening Service

A URL Shortening service built in go using the go-gin framework and MongoDB as the database.

## Getting Started

These instructions will get you a copy of the project up and running.

### Prerequisites

- Docker

### Setup

1. Clone the repo:
```bash
git clone https://github.com/mhope-2/url-shortner.git
```

2. Rename `.env.sample` to `.env` and update env variables accordingly
```bash
cp .env.sample .env
```

3. Run with docker

```bash
docker-compose up
```
or
```
make up
```

## Running tests
```bash
go test -v -cover ./...
```
or
```bash
make test
```

## API  
* Endpoint: `/short-link`  
* Sample Requests
```json
{
  "url": "https://gobyexample.com/random-numbers",
  "slug": "NTMyODk0"
}
```
```json
{
  "url": "https://gobyexample.com/random-numbers"
}
```
* Sample Response  
```json
{
    "result": {
        "shortened_url": "http://localhost:8085/NTMyODk0"
    }
}
```

## Extras
* MongoDB setup guide: https://www.mongodb.com/basics/create-database


