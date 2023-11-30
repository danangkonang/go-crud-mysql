# GOLANG CRUD MYSQL

Simple CRUD app with mysql

# Install

**Run app with docker**

```bash
git clone https://github.com/danangkonang/go-crud-mysql.git

cd go-crud-mysql

docker-compose up --build -d
```

**Run app on local**

`setup your database enviroment in .env file`
```
git clone https://github.com/danangkonang/go-crud-mysql.git

cd go-crud-mysql

go run main.go

```

# Postman

- Download [crud.postman_collection.json](https://raw.githubusercontent.com/danangkonang/nest-micro/master/crud.postman_collection.json)

# Sql

`sql table auto generate after app run`

```sql
CREATE TABLE IF NOT EXISTS users (
  user_id int NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  phone varchar(255) NOT NULL,
  PRIMARY KEY (user_id)
) engine=InnoDB charset=UTF8;
```
