# Name Enricher
The service is a Go-based application that predicts the age, gender, and nationality of a person based on their first name.
The app is integrated with external services [Agify][agify], [Genderize][genderize] and [Nationalize][nationalize].
The service uses a PostgreSQL database for saving enriched user details and also supports updating and deleting user data.

[agify]: https://agify.io/
[genderize]: https://genderize.io/
[nationalize]: https://nationalize.io/

## Features
- makefile
- linter
- dockerfile
- swagger
- logging
- tests (coverage 80%)
- metrics

## Quick start
Currently, wallets-service requires [Go][go] version 1.21 or greater.

[go]: https://go.dev/doc/install

#### Installation guidelines:
```shell
# Clone a repo
$ git clone https://github.com/AlexZav1327/name-enricher.git
# Add missing dependencies
$ go mod tidy
# Start docker container to launch a PostgresSQL database
$ make up
# Run server
$ make run
```
#### Integration tests:
```shell
# App, database and migration
$ make test
```
#### Linters:
```shell
# Run linters
$ make lint
```
## API methods description
### Enrich name
```shell
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "Liza", "surname": "Duchess", "patronymic": "Devonshire"}' \
  'http://localhost:8082/api/v1/user/enrich'
```
#### Response
```json
{"name":"Liza","surname":"Duchess","patronymic":"Devonshire","age":47,"gender":"female","country":"PH"}
```
### Get list of users
```shell
curl -X GET \
  'http://localhost:8082/api/v1/users?textFilter=female&itemsPerPage=2&offset=1&sorting=name&descending=true'
```
#### Response
```json
[
  {"name":"Kitty","surname":"Cat","patronymic":"","age":16,"gender":"female","country":"HK"},
  {"name":"Katherine","surname":"Kit","patronymic":"","age":31,"gender":"female","country":"CL"}
]
```
### Update user
```shell
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{
    "age": 32, 
    "gender": "female",
    "country":"FR"
  }' \
  'http://localhost:8082/api/v1/user/update/Katherine'
```
#### Response
```json
{"name":"Katherine","surname":"Kit","patronymic":"","age":32,"gender":"female","country":"FR"}
```
### Delete wallet
```shell
curl -X DELETE \
'http://localhost:8082/api/v1/user/delete/Katharine'
```