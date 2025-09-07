Federico Bau -- Assignment Results
==================================


### Branches

* **Final Assignment Branch:** https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/federico-bau-assignment
    * [solution/assignment-1-implementation](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1-implementation):
      Main branch with the final assignment implementation implementing what's been requested in [ASSIGNMENT.md](./ASSIGNMENT.md)
        * [Assignment 1-A Refactor Product Repository #1](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/pull/1)
     

### 1) Assignment 1-A: Refactor Product Repository  

* PR: https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/pull/1
* Branch: [solution/assignment-1/a-refactor-product-repository](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1/a-refactor-product-repository)
* Task: First Task of [ASSIGNMENT.md](./ASSIGNMENT.md)'s Catalog endpoint

#### Change Log

* Added `.gitattribute ` to fix EOF
* Created a new Gorm Client wrapper as a struct, defining common methods such as shutdown
* Added `POSTGRES_HOST` + `POSTGRES_SSL_MODE` + `POSTGRES_LOG_LEVEL` to env variables!
* Modified how we load env variables For the database connection, using a `DatabaseConfig` struct 
* Created `DependencyInjectionContainer` which is responsible to create, provide and cleanup applications' resources
  such as database connection. Now `DependencyInjectionContainer` creates the postgres client and the database connection.
* Created `ProductsRepository` as interface and renamed concrete implementation to `SQLProductsRepository`
* `CatalogHandler` now depends on the `ProductsRepository` interface
* `SQLProductsRepository` now uses the `database.Client` struct rather than gorm directly.

### 2) Assignment 1-B: Create a Product Category Model

* PR: ...
* Branch: [solution/assignment-1/b-new-product-category-model](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1/b-new-product-category-model)
* Task: Second Task of [ASSIGNMENT.md](./ASSIGNMENT.md)'s Catalog endpoint, new to create new Product Category Model

#### Change Log

* ...
