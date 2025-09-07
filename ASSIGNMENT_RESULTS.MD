Federico Bau -- Assignment Results
==================================


### Changes from Assignment 1-A Refactor Product Repository  

First Task of [ASSIGNMENT.md](./ASSIGNMENT.md)'s Catalog endpoint


* PR: ...
* [solution/assignment-1/a-refactor-product-repository](https://github.com/Koubae/federico-bau-go-hiring-challenge-1.2.0/tree/solution/assignment-1/a-refactor-product-repository)
  First Task of [ASSIGNMENT.md](./ASSIGNMENT.md)'s Catalog endpoint

#### Change Log

* Added `.gitattribute ` to fix EOF
* Created a new Gorm Client wrapper as a struct, defining common methods such as shutdown
* Added `POSTGRES_HOST` + `POSTGRES_SSL_MODE` to env variables!
* Modified how we load env variables For the database connection, using a `DatabaseConfig` struct 
