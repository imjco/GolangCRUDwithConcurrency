
[![Build Status](https://travis-ci.org/your-username/your-repo.svg?branch=master)](https://travis-ci.org/your-username/your-repo)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Description

An Application developed in Golang. It provides core CRUD (Create, Read, Update, Delete) functionalities in memory, allowing efficient management of book records. The application leverages Golang's concurrency for optimized performance.
## Features
- CRUD Operations: Enables users to perform essential book management tasks: Create, Read, Update, and Delete.

- Concurrency: Utilizes Golang's concurrency to enhance performance during simultaneous operations.

- Unit Testing: Comprehensive unit tests validate the application's functionalities for stability and reliability.
  
- HTML Frontend: Offers an HTML interface for user interaction with the book database.

## Requirements
- Go 1.20 or above  
- Web Browser
- Internet
- VSCode with go extensions installed

## Repository

Clone the repository from github

```bash
  git clone https://github.com/imjco/GolangCRUDwithConcurrency.git

```
## Configurations

### BackEnd
- You can use the steps below to change the Api Url.
  
1. Navigate to the root folder and find "main.go"
2. Scroll to line number 50, as defaut it is set to localhost and port 8080
```go
      fmt.Println("Server started on :8080")
```
- Sample edit. 
```go
     fmt.Println("Server started on 192.168.1.1:3030)
```
    3. Scroll to line number 53, as defaut it is set to localhost and port 8080
```go
      go http.ListenAndServe(":8080", handler.CorsMiddleware(r))
```   
- Sample edit.
```go
      go http.ListenAndServe("192.168.1.1:3030"), handler.CorsMiddleware(r))
```     
    4 Save the file to take effect.

### FrontEnd
- You can use the steps below to change the Api Url.
  
1. Navigate to the project directory and find the "CompanyXYZ_Front" folder.
2. Edit the following JS file.
```js
      APIURL.js
```  
    3. By defaut the url path is set to localhost and port 8080
```js
      http://localhost:8080
```     
    4. Change the url if needed'
    5. Save the file to take effect.

## Run Locally

### Backend
    1. Go to the project directory

```bash
  cd Montani
```

    2. Install dependencies

```bash
  go mod tidy
```

    3. Start the server
  - Note: This will take a few minutes while its Generating 100k+ Random Data

```bash
  go run main.go
```

### Fontend
    1. Go to the project directory and find "CompanyXYZ_Front" folder
    2. Open the file "index.html".
    3. This will open on your default browser.
## Running Tests

     1. Go to the project directory

```bash
  cd GolangCRUDwithConcurrency
```
     2. To test, run the following command.

```bash
  go test -v ./handler
```

```bash
  go test -v ./service
```
