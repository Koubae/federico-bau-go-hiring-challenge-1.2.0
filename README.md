# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

* Repository: https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0
* **Final Assignment Branch:** https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/federico-bau-assignment
* Author: [Federico Bau](https://federicobau.dev/)

### See  Please SEE 🚨🚨🚨

* **[ASSIGNMENT_RESULTS.MD](./dev/ASSIGNMENT_RESULTS.MD)** 
* [PostMan Collection](./dev/mytheresa_(Products)_V001.postman_collection.json)


### QuickStart

#### Hot-Reloader: [Air](https://github.com/air-verse/air) 

* 1) Install [air-verse/air](https://github.com/air-verse/air) globally

```bash
go install github.com/air-verse/air@latest

# Make sure that GOPATH and GOROOT is in your PATH
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin
```

#### Run the application 


```bash
make tidy
# Create Database
make seed 
# Run Server (Hot-Reloader)
make run-reload
# Run Server (No Hot-Reloader)
make run
```

#### Try the API

* Catalog List: http://localhost:8484/catalog?category=Accessories&priceLessThen=10&limit=100&offset=0
* Product Details: http://localhost:8484/catalog/PROD007

I suggest you to use [PostMan Collection](./dev/mytheresa_(Products)_V001.postman_collection.json)

----

## Project Structure

1. **cmd/**: Contains the main application and seed command entry points.

   - `server/main.go`: The main application entry point, serves the REST API.
   - `seed/main.go`: Command to seed the database with initial product data.

2. **app/**: Contains the application logic.
3. **sql/**: Contains a very simple database migration scripts setup.
4. **models/**: Contains the data models and repositories used in the application.
5. `.env`: Environment variables file for configuration.

## Setup Code Repository

1. Create a github/bitbucket/gitlab repository and push all this code as-is.
2. Create a new branch, and provide a pull-request against the main branch with your changes. Instructions to follow.

## Application Setup

- Ensure you have Go installed on your machine.
- Ensure you have Docker installed on your machine.
- Important makefile targets:
  - `make tidy`: will install all dependencies.
  - `make docker-up`: will start the required infrastructure services via docker containers.
  - `make seed`: ⚠️ Will destroy and re-create the database tables.
  - `make test`: Will run the tests.
  - `make run`: Will start the application.
  - `make docker-down`: Will stop the docker containers.

Follow up for the assignemnt here: [ASSIGNMENT.md](ASSIGNMENT.md)
