# Authify Service

The Authify Service handles user authentication and authorization. It communicates with the API Gateway over **gRPC** and uses PostgreSQL for user management.

# Features

- **User Signup**: Register new users with a hashed password.
- **User Login**: Authenticate users with bcrypt password verification and JWT Token generation.
- **Fetch Users**: Checks that user is an **admin**. This endpoint performs cursor pagination to retrieve users from the database.
- **Update User**: Users can update their credentials and profile information.
