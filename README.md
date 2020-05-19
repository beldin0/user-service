Users Microservice
(written in Go)

Start the service with `docker-compose up`. This will create a Postgres database container and the service container.

OpenAPI documentation is in /src/swagger/user

Integration tests can be run with `go test .` This requires ports 8080 and 9000 to be free, and will start a Postgres Docker container to use for the tests.

Criteria:
- Endpoint documentation is auto-generated from proto definitions (in src/proto/user)
- names are stored twice (as-entered and in lowercase) to enable faster text searching of those fields. 

Assumptions:
- Passwords will be encrypted/hashed remotely before being sent to be stored in the database.
- Service will be part of a scaled/managed deployment, so loss of connection with the database will be managed by that (e.g. pod restarting on connection loss)
