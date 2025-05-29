# Go Backend E-commerce Project

## Overview
This project is a backend e-commerce application built using Go. It provides a RESTful API for managing users, products, orders, and addresses. The application is designed to be modular, with separate layers for controllers, services, repositories, and middleware.

## Project Structure
```
Go_backend_ecommerce
├── backend-api
│   ├── cmd
│   │   └── server
│   │       └── main.go
│   ├── internal
│   │   ├── controller
│   │   ├── middleware
│   │   ├── models
│   │   ├── repo
│   │   ├── routers
│   │   └── service
│   ├── migrations
│   ├── pkg
│   ├── global
│   ├── settings
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── go.mod
│   ├── go.sum
│   └── air.toml
├── deployments
│   ├── docker
│   │   ├── postgres
│   │   │   ├── Dockerfile
│   │   │   └── init.sql
│   │   └── nginx
│   │       ├── Dockerfile
│   │       └── nginx.conf
│   └── k8s
│       ├── deployment.yaml
│       ├── service.yaml
│       └── configmap.yaml
├── scripts
│   ├── build.sh
│   ├── deploy.sh
│   ├── migrate.sh
│   └── start-dev.sh
├── docker-compose.yml
├── docker-compose.dev.yml
├── docker-compose.prod.yml
├── .env
├── .env.example
├── .dockerignore
├── Makefile
└── README.md
```

## Getting Started

### Prerequisites
- Go (version 1.16 or higher)
- Docker
- Docker Compose

### Installation
1. Clone the repository:
   ```
   git clone https://github.com/yourusername/Go_backend_ecommerce.git
   cd Go_backend_ecommerce
   ```

2. Navigate to the `backend-api` directory and install dependencies:
   ```
   cd backend-api
   go mod tidy
   ```

### Running the Application
You can run the application using Docker Compose. This will start the Go application along with the PostgreSQL database and Nginx server.

1. Build and start the services:
   ```
   docker-compose up --build
   ```

2. Access the API at `http://localhost:8080/api/v1`.

### Running Migrations
To run database migrations, you can use the provided script:
```
./scripts/migrate.sh
```

### Development
For development, you can use the `air` tool for live reloading. Start the development server with:
```
./scripts/start-dev.sh
```

### Deployment
For deployment, you can use the provided Kubernetes configurations or Docker Compose files for production.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or features.

## License
This project is licensed under the MIT License. See the LICENSE file for details.