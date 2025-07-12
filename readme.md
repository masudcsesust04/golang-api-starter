# Golang JWT Authentication Application

This project provides a JWT authentication system using Golang, with a PostgreSQL backend. It includes user registration, login, logout, and token refresh functionality, as well as CRUD operations for user management.

## Project Setup

- Go 1.23.2 or higher
- A PostgreSQL database

### Installation

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd golang-jwt-auth
    ```

2.  **Install dependencies:**
    ```bash
    go mod download
    ```

3.  **Database Setup:**
    - Create a database in your PostgreSQL server.
    - Run the SQL schema file `schema.sql` to create the necessary tables and indices.

4.  **Environment Variables:**
    - Create a `.env` file in the root of the project.
    - Set the `DATABASE_URL` and `JWT_SECRET` environment variables in the `.env` file:
      ```
      DATABASE_URL=postgres://user_name:password@127.0.0.1:5432/database_name?sslmode=disable
      JWT_SECRET=my-top-secret-key
      ```
    - You can generate a secure secret key with the following command:
      ```bash
      openssl rand -hex 64
      ```

### Running the Server

To start the server, run:

```bash
go run cmd/server/main.go
```

The server will start on port `8080`.

### Running Tests

1.  Create a separate PostgreSQL database for testing (e.g., `ewallet_test`).
2.  Run the `schema.sql` file to create the necessary tables.
3.  Set the `TEST_DATABASE_URL` environment variable:
    ```bash
    export DATABASE_URL=postgres://user_name:password@127.0.0.1:5432/ewallet_test?sslmode=disable
    ```
4.  Run the tests:
    ```bash
    go test ./...
    ```

---

## API Documentation

### Calling Mechanism

-   **Authenticated Routes:** For endpoints that require authentication, you must include an `Authorization` header in your request with a valid JWT access token:
    ```
    Authorization: Bearer <your_access_token>
    ```
-   **Request/Response Format:** All request and response bodies are in JSON format. Ensure your requests have the `Content-Type: application/json` header.

---

### Authentication APIs

These endpoints handle user login, logout, and token management.

#### 1. Login

-   **Description:** Authenticates a user and returns an access token and a refresh token.
-   **Method:** `POST`
-   **Path:** `/login`
-   **Authentication:** Not required.
-   **Request Body:**
    ```json
    {
      "email": "user@example.com",
      "password": "your_password"
    }
    ```
-   **Success Response (200 OK):**
    ```json
    {
      "access_token": "...",
      "refresh_token": "..."
    }
    ```

#### 2. Refresh Access Token

-   **Description:** Generates a new access token using a valid refresh token.
-   **Method:** `POST`
-   **Path:** `/token/refresh`
-   **Authentication:** Not required.
-   **Request Body:**
    ```json
    {
      "user_id": 1,
      "refresh_token": "your_refresh_token"
    }
    ```
-   **Success Response (200 OK):**
    ```json
    {
      "access_token": "..."
    }
    ```

#### 3. Logout

-   **Description:** Invalidates the user's refresh token on the server. The client is responsible for deleting the tokens.
-   **Method:** `POST`
-   **Path:** `/logout`
-   **Authentication:** **Required**.
-   **Request Body:**
    ```json
    {
      "user_id": 1
    }
    ```
-   **Success Response:** `204 No Content`

---

### User Management APIs

These endpoints handle CRUD operations for users.

#### 1. Create User

-   **Description:** Creates a new user account.
-   **Method:** `POST`
-   **Path:** `/users`
-   **Authentication:** Not required.
-   **Request Body:**
    ```json
    {
      "first_name": "John",
      "last_name": "Doe",
      "phone_number": "1234567890",
      "email": "john.doe@example.com",
      "password": "a_strong_password",
      "status": "active"
    }
    ```
-   **Success Response (201 Created):** The newly created user object.

#### 2. Get All Users

-   **Description:** Retrieves a list of all users.
-   **Method:** `GET`
-   **Path:** `/users`
-   **Authentication:** **Required**.
-   **Success Response (200 OK):** An array of user objects.

#### 3. Get User by ID

-   **Description:** Retrieves a single user by their ID.
-   **Method:** `GET`
-   **Path:** `/users/{id}` (e.g., `/users/1`)
-   **Authentication:** **Required**.
-   **Success Response (200 OK):** A single user object.

#### 4. Update User

-   **Description:** Updates an existing user's information.
-   **Method:** `PUT`
-   **Path:** `/users/{id}` (e.g., `/users/1`)
-   **Authentication:** **Required**.
-   **Request Body:**
    ```json
    {
      "first_name": "John",
      "last_name": "Doe",
      "phone_number": "1234567890",
      "email": "john.doe@example.com",
      "status": "active"
    }
    ```
-   **Success Response (200 OK):** The updated user object.

#### 5. Delete User

-   **Description:** Deletes a user by their ID.
-   **Method:** `DELETE`
-   **Path:** `/users/{id}` (e.g., `/users/1`)
-   **Authentication:** **Required**.
-   **Success Response:** `204 No Content`