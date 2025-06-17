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
