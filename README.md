# Event manager API

This API solution is part of b2match's application task for job interview. It is a simple REST API for managing users, events and event meetings.
It is written in Golang and uses Gin framework for routing and Gorm for ORM.

## Installation
To install this API locally, follow these steps:
1. **Clone the repository**: Begin by cloning the repository to your local machine. You can use the following command:
    ```bash
    git clone github.com/KarloVrd/event-manager-API
2. **Navigate to the project directory**: Change your current working directory to the root of the cloned repository:
    ```bash
    cd event-manager-API
    ```
3. **Install dependencies**: Install the necessary dependencies for your API. You can use the following command:
    ```bash
    go mod download
    ```
## Usage
You can run server and start sending it requests at http://127.0.0.1:8080 by running this command:
```bash
go run main.go
```

You can also run already written test at main_test.go:
```bash
go test
```