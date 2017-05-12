# API - *peragrin api*

## Development

### Required Software

* [docker](https://docs.docker.com/)
* [git](https://git-scm.com/)

### Getting Started

1. `git clone https://gitlab.com/peragrin/api.git $GOPATH/src/gitlab.com/peragrin/api`

2. `docker-compose up -d db`

3. `cd $GOPATH/src/gitlab.com/peragrin/api`

4. `go run main.go migrate`

5. `go run main.go createfixturedata`

6. `curl -u jteppinette:jteppinette localhost:8000/account`

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
* LOCATIONIQ_API_KEY

### Docker

1. `docker build . -t app`

2. `docker run \
      -d
      -e POSTGRES_DB=db
      -e POSTGRES_USER=db
      -e POSTGRES_PASSWORD=db-secret
      --name db
      postgres:9.6.2`

3. `docker run
      -d
      -p 8000:80
      -e DB_NAME=db
      -e DB_USER=db
      -e DB_PORT=3306
      -e DB_PASSWORD=db-secret
      -e DB_HOST=db
      --link db
      --name app
      app`

4. `docker exec -it app api migrate`

5. `docker exec -it app api createfixturedata`
