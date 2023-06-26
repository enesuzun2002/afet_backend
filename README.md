# Disaster Management App Backend

This is the backend code for a disaster management app. The backend is built using GoLang and utilizes various libraries and frameworks. The app allows users to sign up, log in, and submit their location in case of a disaster. It also provides endpoints to retrieve information about users in need of help.

## Features

- User sign up
- User login
- Submitting user location during a disaster
- Retrieving information about users in need of help

## Installation

To run the backend locally, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/enesuzun2002/afet_backend.git
   ```

2. Change into the project directory:

   ```bash
   cd afet_backend
   ```

3. Install the dependencies:

   ```bash
   go mod download
   ```

4. Set up the PostgreSQL database and update the connection string in the code.

5. Run the project:

   ```bash
   go run .
   ```

6. The backend server will start running on `http://localhost:8080`.

## API Endpoints

The following API endpoints are available:

- **POST /signup**: Creates a new user account.
- **POST /login**: Logs in an existing user.
- **POST /yardimaihtiyacimvar**: Submits user location during a disaster (requires authentication).
- **GET /yardimedebilirim**: Retrieves information about users in need of help (requires authentication).

## Authentication

The backend uses JSON Web Tokens (JWT) for authentication. When a user signs up or logs in, a JWT token is generated and returned in the response. This token should be included in the `Authorization` header of subsequent requests as follows:

```
Authorization: <token>
```

## Database Setup

The backend uses a PostgreSQL database. Make sure you have PostgreSQL installed and create a new database. Update the connection string in the code (`dsn` variable) to match your database configuration.

## Dependencies

The backend relies on the following external libraries:

- [github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go): Used for JWT token generation and validation.
- [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin): HTTP web framework for building the API.
- [gorm.io/driver/postgres](https://gorm.io/driver/postgres): PostgreSQL database driver for GORM.
- [gorm.io/gorm](https://gorm.io/gorm): ORM library for database operations.

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
