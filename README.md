# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

* Repository: https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0
* **Final Assignment Branch:** https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/federico-bau-assignment
* Author: [Federico Bau](https://federicobau.dev/)

### Branches

* **Final Assignment Branch:** https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/federico-bau-assignment
  * [solution/assignment-1-implementation](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1-implementation): 
    Main branch with the final assignment implementation implementing what's been requested in [ASSIGNMENT.md](./ASSIGNMENT.md)
    * [solution/assignment-1/a-refactor-product-repository](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1/a-refactor-product-repository)
      First Task of [ASSIGNMENT.md](./ASSIGNMENT.md)'s Catalog endpoint


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
