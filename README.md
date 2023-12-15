# micro-fiber-test

postgreSQL container: scripts/docker_postgres.txt

redis container: scripts/docker-redis.yml

Redis is used as session storage backend

database schema: migrations/20220905_usm_init.up.sql

REST endpoints: scripts/Insomnia.json ==> [Insomina](https://insomnia.rest/download)

Github authentication (OAuth2 integration) ==> [Homepage](https://localhost:8443/index.html)

Update all dependencies: go get -u then go mod tidy

Build:

    go build -ldflags "-s -w" -o micro-fiber-test.exe

-w: remove debugging information

-s: remove symbol table

TIPS:

- List possible platforms: go tool dist list
- Operating System & Architecture: go env GOOS GOARCH
- Build for target operating system and architecture: GOOS=linux GOARCH=ppc64 go build
- Assign new_value to variable variable_name in package_path:
  - go build -ldflags="-X 'package_path.variable_name=new_value'
  - go build -ldflags="-X 'main.Version=v1.0.0'"
- CGO (cross-compile native support for target ploatform):
  - https://stackoverflow.com/questions/61515186/when-using-cgo-enabled-is-must-and-what-happens
  - https://stackoverflow.com/questions/64531437/why-is-cgo-enabled-1-default
- List ldflags: go build --ldflags="--help"
- List dependencies updates: go list -m -u all
- Update all dependencies: go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
- Go installation on Debian: https://www.digitalocean.com/community/tutorials/how-to-install-go-on-debian-10
- Generating a self signed certificate: run cmd/certSelfSigned
- Prometheus metrics exposed by default on "/metrics" path

Configuration: config/config.yaml

- PostgreSQL:
  - pgUrl: connection url (e.g: postgres://${user}:${password}@${host}:${5433}/${database})
  - pgPoolMin: Connection pool min size
  - pgPoolMax: Connection pool max size
- accessLogFile: Access log file (access.log)
- stdLogFile: Standard log file (micro-fiber-test.log)
- OAuth2 - Gitlab:
  - oauthClientId: clientId for github connection
  - oauthClientSecret: client secret for github connection
  - oauthGithub: github authorize url
  - oauthCallback: github access token url
  - oauthRedirectUri: redirect url on local app
  - oauthDebug: Enable/disable debug logs
- Redis:
  - redisHost: Redis instance host
  - redisPort: Redis port (defaults to 6379)
  - redisUser: Redis account username (Default blank)
  - redisPass: Redis account password (Default blank)
- prometheusEnabled: Enable/Disable prometheus middleware
- metricsPath: Prometheus exposition path (Defaults to "/metrics")
