# API using Go Standard Library

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Endpoints](#endpoints)
- [Why Go's Standard Library?](#why-gos-standard-library)
- [Contributing](#contributing)

## Introduction

This project is a RESTful API built with Go, utilizing only its standard library. It serves as a basic CRUD (Create, Read, Update, Delete) application for managing user data.

Tests

## Features

- Rate limiting
- Logging using `slog`
- Data validation
- Error handling

## Requirements

- Go 1.21

## Installation

1. Clone this repository.

   ```bash
   git clone https://github.com/your-username/api-std-lib.git
   ```

2. Install dependencies (postgres driver).

   ```bash
    go get
   ```

3. Run the application.

   ```bash
    go run main.go

   ```

## Endpoints

- `POST /user` : Create a new user
- `GET /user` : Fetch all users
- `GET /user/{id}` : Fetch a single user by ID
- `PUT /user/{id}` : Update a user by ID
- `DELETE /user/{id}` : Delete a user by ID

## Why Go's Standard Library?

The purpose of using Go's standard library for this project is to demonstrate that it's possible to build a highly efficient, fast, and secure RESTful API without relying on third-party libraries or frameworks. This approach leads to lightweight, easy-to-maintain code and allows for greater control over the application's behavior. Key benefits include:

- Simplicity: Reduced code complexity by utilizing native functionalities.
- Performance: Highly optimized, native solutions for common tasks like HTTP routing, data encoding/decoding, and database connection pooling.
- Security: Reduced surface area for security vulnerabilities by relying on well-tested, standard implementations.
- Learning: Great way to understand the inner workings of various components like caching, logging, and rate limiting by implementing them from scratch.

## Contributing

If you'd like to contribute, please fork the repository and create a new branch, then submit a pull request.
