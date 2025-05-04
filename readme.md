# Golang JWT authentication application

## Application overview
This project provides a JWT authentication using golang.

## Project Setup
- Go 1.23.2 or higher 
- A PostgreSQL database

### Installation
1. Clone the code repository:
```bash
git clone <repository-url>
cd golang-jwt-auth
```

2. Install dependencies:
```bash
go mod download
```

3. Database setup
- Create a database in your PostgreSQL server.
- Run the SQL schema file `schem.sql` to create necessary table and indices.

4. Set `DATABASE_URL` to environment variable:
```bash
export DATABASE_URL=postgres://user_name:password@127.0.0.1:5432/database_name?sslmode=disable
```
5. Setup the `JWT_SECRET` environment variable:
```bash
   export JWT_SECRET="my-top-secret-key"
```

Note: Generate secret key
```bash
openssl rand -hex 64
```

## Running the server:
To start the server, run:
```bash
go run cmd/server/main.go
```

The server will start on port `8080`


## Run test:
1. Create a postgresql test database `ewallet_test` and create tables defined in `schema.sql` file.
2. Set `TEST_DATABASE_URL` to environment variable:
```bash
export DATABASE_URL=postgres://user_name:password@127.0.0.1:5432/ewallet_test?sslmode=disable
``` 
3. Run test:
```bash
go test ./...
```

API End point testing:
1. Create user:
```bash
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{"first_name": "md", "last_name": "rana", "phone_number": "809890899", "email": "rana@gmail.com", "password": "example123", "status": "active"}'
```

1. Login
```bash
curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{"email": "rana@gmail.com", "password": "example123"}'
```

2. Refresh token:
```bash
curl -X POST http://localhost:8080/token/refresh \
-H "Content-Type: application/json" \
-d '{"refresh_token": ""}'
```

3. Logout:
```bash
curl -X POST http://localhost:8080/logout \
-H "Content-Type: application/json" \
-d '{"refresh_token": "refresh_token_here"}'
```

4. List of user:
```bash
curl -X GET http://localhost:8080/users \
-H "Authorization: Bearer jwt_access_token_here"
```

5. Get user by id:
```bash
curl -X GET http://localhost:8080/users/1 \
-H "Authorization: Bearer jwt_access_token_here"
```

6. Update user:
```bash
curl -X PUT http://localhost:8080/users/1 \
-H "Authorization: Bearer jwt_access_token_here" \
-d '{"phone_number": "8098908080"}'
```

7. Delete id:
```bash
curl -X DELETE http://localhost:8080/users/1 \
-H "Authorization: Bearer jwt_access_token_here"
```
