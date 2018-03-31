# Peragrin API - *peragrin api*

## Development

### Required Software

* [docker](https://docs.docker.com/)
* [git](https://git-scm.com/)

### Getting Started

1. `git clone https://github.com/jteppinette/peragrin-api.git $GOPATH/src/github.com/jteppinette/peragrin-api`

2. `docker-compose up -d db minio mail`

3. `cd $GOPATH/src/github.com/jteppinette/peragrin-api`

4. `go run main.go migrate -m <migrations-directory>`

5. `go run main.go serve`

## Usage

### Environment Variables

Any variables marked as `insecure: true` should be overriden before being added to a production system.

* DB_NAME             `default: db`
* DB_USER             `default: db`
* DB_PASSWORD         `defualt: secret, insecure: true`
* DB_HOST             `default: 0.0.0.0`
* DB_PORT             `default: 5432`
* PORT                `default: 8000`
* TOKEN_SECRET        `default: token-secret, insecure: true`
* LOG_LEVEL           `default: info`
* LOCATIONIQ_API_KEY  `insecure: true`
* MANDRILL_KEY        `insecure: true`
* MAIL_FROM           `default: notifications@peragrin.localhost`
* MAIL_HOST           `default: 0.0.0.0`
* MAIL_PORT           `default: 1025`
* MAIL_PASSWORD       `insecure: true`
* MAIL_USER           `insecure: true`
