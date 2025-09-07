Federico Bau -- Assignment Results
==================================


### Branches

* **Final Assignment Branch:** https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/federico-bau-assignment
    * [solution/assignment-1-implementation](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1-implementation):
      Main branch with the final assignment implementation implementing what's been requested in [ASSIGNMENT.md](./ASSIGNMENT.md)
        * [1) Assignment 1-A Refactor Product Repository #1](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/pull/1)
        * [2) Assignment 1-B: Create a Product Category Model #2](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/pull/2)
     

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

* PR: https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/pull/2
* Branch: [solution/assignment-1/b-new-product-category-model](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1/b-new-product-category-model)
* Task: Second Task of [ASSIGNMENT.md](./ASSIGNMENT.md)'s Catalog endpoint, new to create new Product Category Model

#### Change Log

* Change multiple aspects of the database "seed" and migration. Most likely, there should be in place a proper migration tool.
  If creating SQL file is a requirements it could be done in 2 steps where it automatically loads the curren
  DDL schema by the GORM models (without write them manually especially becase ORM and query will have discrepancies)
  For now what this does is: 1) Drop all models 2) Create All models 3) Load the seed data. So we let GORM's ORM do its job

* Now [sql](./sql) folder contains only SEED / Test data, but wont' create nor delete tables.
* Refactor models into app/models
* Refactor SQLProductsRepository into repository package
* Added Category in Product in catalog response

### 3) Assignment 1-C: Add Pagination in Catalog Endpoint

* PR: https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/pull/3
* Branch:  [solution/assignment-1/c-add-pagination-in-catalog-endpoint](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1/c-add-pagination-in-catalog-endpoint)
* Task: Third Task of [ASSIGNMENT.md](./ASSIGNMENT.md)'s Catalog endpoint, Update the catalog handler and relevant repositories to support offset pagination.

#### Change Log

* Created `PaginationRequest` struct that validates and loads limit and offset from request
* Added pagination to handler and repository
* Fixed Api tests by fixing `OKResponse` and `ErrorResponse`. also added them in Handler
* Added `RecoverPanic` and `LogAccessMiddleware` middlewares

### 4) Assignment 1-D: Catalog Endpoint query by category and price

* PR: https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/pull/4
* Branch: [solution/assignment-1/d-catalog-endpoint-query-by-category-and-price](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1/d-catalog-endpoint-query-by-category-and-price)
* Task: Fourth Task of [ASSIGNMENT.md](./ASSIGNMENT.md)'s Catalog endpoint, Catalog Endpoint query by category and price

#### Change Log

* Added category and price query to handler and repository


### 5) Assignment 1-E: Product details endpoint

* PR: https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/pull/5
* Branch: [solution/assignment-1/e-product-detail-endpoint](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1/e-product-detail-endpoint)
* Task: Fifth Task of [ASSIGNMENT.md](./ASSIGNMENT.md)'s Product details endpoint

#### Change Log

* Renamed Catalog Handler `HandleGet` --> `ListCatalog`
* Created new endpoint `/catalog/:code` at handler `CatalogHandler.GetProductDetails`
* Implemented `ProductDetailsHandler`

### Example

Request: 

```bash
curl --location 'http://localhost:8484/catalog/PROD007'
```

Response: 
```json
{
  "id": 3,
  "code": "PROD007",
  "price": 18.2,
  "category": {
    "code": "CAT001",
    "name": "Clothing"
  },
  "variants": [
    {
      "id": 17,
      "sku": "SKU007A",
      "name": "Variant A",
      "price": 18.2
    },
    {
      "id": 18,
      "sku": "SKU007B",
      "name": "Variant B",
      "price": 18.2
    },
    {
      "id": 19,
      "sku": "SKU007C",
      "name": "Variant C",
      "price": 18.2
    },
    {
      "id": 20,
      "sku": "SKU007D",
      "name": "Variant D",
      "price": 18.2
    },
    {
      "id": 21,
      "sku": "SKU007E",
      "name": "Variant E",
      "price": 18.75
    }
  ]
}
```


### 6) Assignment 1-F: Categories endpoint CRUD

* PR: 
* Branch: [solution/assignment-1/f-category-crud]()
* Task: Sixth Task of [ASSIGNMENT.md](./ASSIGNMENT.md)'s Categories endpoint CRUD

#### Change Log

* Moved structured used by handlers into `dto` package
* Implemented List Category Endpoint


#### List categories

Request:

```bash
curl --location 'http://localhost:8484/categories?limit=10&offset=0'
```

Response:

```bash
{
    "total": 3,
    "categories": [
        {
            "code": "CAT001",
            "name": "Clothing"
        },
        {
            "code": "CAT002",
            "name": "Shoes"
        },
        {
            "code": "CAT003",
            "name": "Accessories"
        }
    ]
}
```