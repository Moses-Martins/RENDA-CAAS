# Renda-CAAS 
> A simple backend service in Go with MongoDB, OAuth2 authentication, and role-based access control for multiple products.

## Table of Contents
- [Project Overview](#project-overview)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Environment Variables](#environment-variables)
- [Endpoints](#endpoints)
- [Examples](#examples)
- [Structure](Structure.md)
- [License](LICENSE)
## Project Overview
This backend service is designed as a central authentication and authorization system for Renda’s products: Renda360, SCALE, and Project Horizon.
All products share a single user database, but each product can be managed and accessed independently.

### Registration
- Three separate registration endpoints (`/v1/register/renda360`, `/v1/register/scale`, `/v1/register/horizon`) allow each product to have its own onboarding flow.
This was built under the assumption that each product could operate as a standalone service, but all user data is unified in one MongoDB database.

- Regardless of which endpoint is used, a new user is created in the central database and assigned a role for each product.

### Login with Email/Password
- If the user registers through any product endpoint, they are created in the database and assigned the `"User"` role for that product and `"Viewer"` for the others.
On login, the user’s existing roles are used and a JWT is issued.

### Login with Google
- If the user is new, they are automatically created in the database and assigned the `"User"` role for all products (`Renda360`, `Scale`, and `Horizon`).
- If the user already exists, their existing roles are used.

### Roles & Product Access
- Role-based access control is enforced per product.
Each user has a productRoles map, e.g.:
```json
{
  "Renda360": "User",
  "Scale": "Viewer",
  "Horizon": "Viewer"
}
```
- Roles supported: `Superadmin`, `Admin`, `User`, `Viewer`.

### Role Permissions

- **Superadmin:**  
  - Can access all endpoints, including admin-only routes.
  - Can assign or remove the `"Admin"` role for any user and product.
  - Can update any user's role for any product, including removing all access. Which means the user has no access to that product at all.

- **Admin:**  
  - Can access product dashboards for their product.
  - Can promote/demote users to `"User"`, `"Viewer"`, or remove access completely within their product.
  - Cannot assign or remove `"Admin"` roles.

- **User:**  
  - Can access the `/v1/dashboard/{product}` endpoint for products where they have the `"User"` role.
  - Can view their own info via `/v1/me`.

- **Viewer:**  
  - Can only view their own info via `/v1/me`.
  - Cannot access product dashboards or admin endpoints.

**Note:**  
- Only users with `"User"`, `"Admin"`, or `"Superadmin"` roles for a product can access `/v1/dashboard/{product}` for that product.
- Only `"Superadmin"` and product `"Admin"` roles can update privileges via `/v1/admin/update-privilege`.

### OAuth2 & JWT
- Users can log in with email/password or Google OAuth2.
- On login, a JWT is issued containing user info and their roles for each product.
- All protected routes check the JWT for product-specific permissions.

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

    # Google OAuth2 Configuration
    GOOGLE_CLIENT_ID=your_google_client_id
    GOOGLE_CLIENT_SECRET=your_google_client_secret
    GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
    ```

    > **Note:**  
    > - To use Google login, you must [register your app with Google](https://console.developers.google.com/) and obtain a Client ID and Client Secret.  
    > - Set the redirect URL in the Google console to match `GOOGLE_REDIRECT_URL` above.

    Start MongoDB: Make sure you have access to a MongoDB instance (local or Atlas).
    
   
## Usage

1. Start the server:
    ```sh
    go run main.go
    ```
2. The server will start on the port specified in the `.env` file (default is `8080`).

## Environment Variables

- `MONGO_URI`: The URI for connecting to the MongoDB database. (Make sure your username and database password is in the link.)
- `PORT`: The port on which the server will run.



## Endpoints

Here is a list of all the endpoints and their respective methods. You can access endpoints using tools like Postman or curl.

| Endpoint | Method | Parameters | Description |
| :------: | :----: | :--------- | :---------- |
| /v1/register/renda360 |	POST | name, email, password, confirm password | Register a user for Renda360 (User becomes a Viewer on other products) |
| /v1/register/scale | POST | name, email, password, confirm password | Register a user for Scale (User becomes a Viewer on other products) |
| /v1/register/horizon | POST | name, email, password, confirm password | Register a user for Horizon (User becomes a Viewer on other products) |
| /v1/login | POST | email, password | Authenticates a user and returns a JWT token |
| /auth/google/login | GET | | Redirects user to Google for OAuth2 login. On success, creates or logs in the user and returns a JWT token |
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
  "password": "securePassword123",
  "confirmPassword": "securePassword123"
}
```

</details> <details> <summary><strong>Register a user (Scale)</strong></summary>

POST `/v1/register/scale`

```json
{
  "name": "John Smith",
  "email": "john@example.com",
  "password": "anotherSecurePassword",
  "confirmPassword": "anotherSecurePassword"
}
```
</details> <details> <summary><strong>Register a user (Horizon)</strong></summary>

POST `/v1/register/horizon`
```json
{
  "name": "Alice Johnson",
  "email": "alice@example.com",
  "password": "MySafePass456",
  "confirmPassword": "MySafePass456"
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

</details> <details> <summary><strong>Login with Google</strong></summary>

GET `/auth/google/login`

**How to use:**
1. Visit `/auth/google/login` in your browser.
2. You will be redirected to Google to log in.
3. After successful login, you will be redirected back and receive a JWT token in the response.

✅ Example response:
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


