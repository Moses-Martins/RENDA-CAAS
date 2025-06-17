# Renda-CAAS 
> A simple backend service in GO with MongoDB and JWT Auth

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Environment Variables](#environment-variables)
- [Examples](#examples)
- [Contributing](CONTRIBUTING.md)
- [License](LICENSE)
## Getting Started

### Prerequisites

Make sure you have GO (>= 1.24.3) installed on your machine.

### Installation

1. Clone the repository:
    ```properties
    git clone https://github.com/Moses-Martins/RENDA-CAAS.git
    cd RENDA-CAAS
    ```
2. Install the dependencies:
    ```properties
    go mod download
    ```

3. Create a `.env` file in the root directory and add the following environment variables:
    ```properties
    # MongoDB Configuration
    MONGODB_URI=mongodb+srv://..............
   
    # Server Configuration
    PORT=8080 
    
    # JWT Secret Key
    JWT_SECRET=your_jwt_secret
    ```
## Usage

1. Start the server:
    ```sh
    go run main.go
    ```
2. The server will start on the port specified in the `.env` file (default is `8080`).

## Environment Variables

- `MONGODB_URI`: The URI for connecting to the MongoDB database.
- `PORT`: The port on which the server will run.



## References

go-purple exposes multiple API endpoints by defult, mainly used for authentication and testing purposes. Here is a list of all the endpoints and their respective methods. You can test these directly by issuing the requests in `tests/*.http` files.

| Endpoint | Method | Parameters | Description |
| :------: | :----: | :--------- | :---------- |
| / | GET | | Static binding for files in `public` |
| /api/v1/auth/login | POST | username, password | Attempts to authenticate a user. If successful, it returns a web token |
| /api/v1/auth/logout | POST | | Invalidates the current session |
| /api/v1/auth/register | POST | username, password, email | Attempts to create a new user. The provided username **must** be unique |
| /api/v1/users | GET | | Returns a list of all users. Authorization is required to access this endpoint |
| /api/v1/users/`:user` | GET | | Returns a specific user. Authorization is required to access this endpoint |




| Endpoint | Method | Parameters | Description |
| :------: | :----: | :--------- | :---------- |
| /renda360/register |	POST | name, email, password | Register a user for Renda360 (user gets Viewer on others) |
| /scale/register | POST | name, email, password | Register a user for Scale (user gets Viewer on others) |
| /horizon/register | POST | name, email, password | Register a user for Horizon (user gets Viewer on others) |
| /login | POST | email, password | Authenticates a user and returns a JWT token |
| /me |	GET | JWT token (header) | Returns user info and product access details |
| /dashboard/{product} | GET | JWT token (header), product name | Product dashboard, only for Admin/User/SuperAdmin of the product |
| /admin/update-privilege | POST | email, product, role, JWT token | Superadmin or product admin updates a user's role for a product |
