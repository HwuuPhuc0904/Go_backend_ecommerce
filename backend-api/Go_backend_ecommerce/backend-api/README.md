# Go Backend E-commerce API

This project is a backend API for an e-commerce application built using Go. It follows a clean architecture pattern, separating concerns into different layers such as controllers, services, repositories, and models.

## Project Structure

```
Go_backend_ecommerce
├── backend-api
│   ├── cmd
│   │   └── server
│   │       └── main.go          # Entry point of the application
│   ├── internal
│   │   ├── controller            # Contains business logic and interacts with models
│   │   ├── middleware            # Middleware functions for request processing
│   │   ├── models                # Data models defining the structure of the data
│   │   ├── repo                  # Repository layer for data access
│   │   ├── routers               # Routing logic for API endpoints
│   │   └── service               # Service layer containing business logic
│   ├── migrations                 # Migration files for database schema management
│   ├── pkg                       # Utility packages for reuse across the application
│   ├── global                    # Global variables and configurations
│   ├── settings                  # Configuration settings for the application
│   ├── Dockerfile                # Dockerfile for building the backend API image
│   ├── .dockerignore             # Files and directories to ignore in Docker builds
│   ├── go.mod                    # Go module file defining dependencies
│   ├── go.sum                    # Checksums for module dependencies
│   └── README.md                 # Documentation for the backend API
├── docker-compose.yml             # Defines services, networks, and volumes for Docker
├── docker-compose.dev.yml         # Development-specific configurations for Docker Compose
├── docker-compose.prod.yml        # Production-specific configurations for Docker Compose
├── .env                           # Environment variables for the application
├── .env.example                   # Example of required environment variables
├── nginx
│   ├── Dockerfile                 # Dockerfile for building the Nginx server image
│   └── nginx.conf                 # Configuration for the Nginx server
├── postgres
│   ├── init.sql                   # SQL commands to initialize the PostgreSQL database
│   └── Dockerfile                 # Dockerfile for building the PostgreSQL image
├── scripts
│   ├── build.sh                   # Script to build Docker images
│   ├── deploy.sh                  # Script to deploy the application
│   └── start-dev.sh               # Script to start the application in development mode
└── README.md                      # Documentation for the overall project
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- Docker and Docker Compose

### Installation

1. Clone the repository:

   ```
   git clone <repository-url>
   cd Go_backend_ecommerce
   ```

2. Build the Docker images:

   ```
   ./scripts/build.sh
   ```

3. Start the application:

   ```
   docker-compose up
   ```

### Usage

- The API will be available at `http://localhost:8080/api/v1`.
- Use tools like Postman or curl to interact with the API endpoints.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

### License

This project is licensed under the MIT License. See the LICENSE file for details.