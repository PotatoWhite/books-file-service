
1. postgresql - docker

```shell
docker run -p 5432:5432 --name books_postgres -e POSTGRES_PASSWORD=1234 -d postgres
```

2. postgresql - account

```shell
create database file_service;
create users fileaccount with encrypted password '1234';
grant all privileges on database file_service to fileaccount;
```

3. go - gqlgen

```shell
go get -u github.com/99designs/gqlgen
go run github.com/99designs/gqlgen init
```

4. go - gqlgen - if you need to update schema

```shell
go run github.com/99designs/gqlgen generate
``` 

5. gqlgen - dynamic resolver 설정
- Folder 

5. go - kafka

```shell
go get github.com/confluentinc/confluent-users-go/users
```