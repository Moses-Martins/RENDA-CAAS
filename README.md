# Renda-CAAS 
> A simple backend service in GO with MongoDB and JWT Auth

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Environment Variables](#environment-variables)
- [Endpoints](#endpoints)
- [Examples](#examples)
- [Structure](Structure.md)
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
    MONGO_URI=mongodb+srv://<UserName>:<db_password>@cluster0.gef1jsp.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
    MONGO_DB=my_database
   
    # Server Configuration
    PORT=8080 
    
    # JWT Secret Key
    JWT_SECRET=your_jwt_secret
    ```
    Start MongoDB: Make sure you have access to a MongoDB instance (local or Atlas).
   
## Usage

1. Start the server:
    ```sh
    go run main.go
    ```
2. The server will start on the port specified in the `.env` file (default is `8080`).

## Environment Variables

- `MONGO_URI`: The URI for connecting to the MongoDB database.
- `PORT`: The port on which the server will run.



## Endpoints

Here is a list of all the endpoints and their respective methods. You can access endpoints using tools like Postman or curl.

| Endpoint | Method | Parameters | Description |
| :------: | :----: | :--------- | :---------- |
| /v1/register/renda360 |	POST | name, email, password | Register a user for Renda360 (User becomes a Viewer on other products) |
| /v1/register/scale | POST | name, email, password | Register a user for Scale (User becomes a Viewer on other products) |
| /v1/register/horizon | POST | name, email, password | Register a user for Horizon (User becomes a Viewer on other products) |
| /v1/login | POST | email, password | Authenticates a user and returns a JWT token |
| /v1/me |	GET | JWT token (header) | Returns user info and product access details |
| /v1/dashboard/{product} | GET | JWT token (header), product name | Product dashboard, only for Admin/User/SuperAdmin of the product |
| /v1/admin/update-privilege | PATCH | email, product, role, JWT token | Superadmin or product admin updates a user's role for a product |


## Examples

<details> <summary><strong>Register a user (Renda360)</strong></summary>

POST `/v1/register/renda360`

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "securePassword123"
}
```

</details> <details> <summary><strong>Register a user (Scale)</strong></summary>

POST `/v1/register/scale`

```json
{
  "name": "John Smith",
  "email": "john@example.com",
  "password": "anotherSecurePassword"
}
```
</details> <details> <summary><strong>Register a user (Horizon)</strong></summary>

POST `/v1/register/horizon`
```json
{
  "name": "Alice Johnson",
  "email": "alice@example.com",
  "password": "MySafePass456"
}
```
</details> <details> <summary><strong>Login</strong></summary>

POST `/v1/login`
```json
{
  "email": "jane@example.com",
  "password": "securePassword123"
}
```
✅ Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
</details> <details> <summary><strong>Get User Info</strong></summary>

GET `/v1/me`

Headers:

`Authorization`: Bearer <JWT_TOKEN>

</details> <details> <summary><strong>Product Dashboard</strong></summary>

GET `/v1/dashboard/renda360` (or scale, horizon)

Headers:

`Authorization`: Bearer <JWT_TOKEN>

</details> <details> <summary><strong>Update User Privilege</strong></summary>

PATCH `/v1/admin/update-privilege`

Headers:

`Authorization`: Bearer <JWT_TOKEN>

Body:
```json
{
  "email": "user@example.com",
  "product": "Scale",
  "role": "Admin"
}
```
Only users with the **Superadmin** role can assign or remove the ***Admin*** role for any product.
Product admins can only promote or demote users to ***User*** or ***Viewer*** roles within their own product, but cannot assign or remove the *Admin* role for any user.

</details>


