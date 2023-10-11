# XYZ Books Web Application
[![Build Status](https://travis-ci.org/your-username/your-repo.svg?branch=master)](https://travis-ci.org/your-username/your-repo)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Description

The "XYZ Books Web Application" is a dynamic and user-friendly online platform designed to cater to all your reading needs. Crafted with a combination of HTML and Golang, this web application offers an immersive experience for book enthusiasts and readers of all ages.

## Features
- Extensive Book Catalog: Explore a vast collection of books spanning various genres, from classic literature to contemporary bestsellers.

- User-Friendly Interface: Navigate effortlessly through our intuitive user interface designed with HTML to discover new books, read reviews, and manage your reading list.

- Secure Transactions: Our Golang backend ensures secure and efficient handling of user data and transactions, guaranteeing a safe environment for all users.

- Responsive Design: Enjoy seamless access to the XYZ Books Web Application from your desktop, tablet, or mobile device.

- Whether you're a seasoned bookworm or just beginning your literary journey, the XYZ Books Web Application promises to be your go-to destination for discovering, exploring, and enjoying the world of books.

Start your reading adventure today with XYZ Books!
## Requirements
- Go 1.20 or above  
- Web Browser
- Internet
- VSCode with go extensions installed

## Repository

Clone the repository from github

```bash
  git clone https://github.com/imjco/Montani.git

```
## Configurations

### BackEnd
    1. Navigate to the root folder and find "main.go"
    2. Scroll to line number 50, as defaut it is set to localhost and port 8080
```bash
      fmt.Println("Server started on :8080")
```
    - change it if needed. eg. 
```bash
     fmt.Println("Server started on 192.168.1.1:3030)
```
    3. Scroll to line number 53, as defaut it is set to localhost and port 8080
```bash
      go http.ListenAndServe(":8080", handler.CorsMiddleware(r))
```   
    - change it if needed. eg.
```bash
      go http.ListenAndServe("192.168.1.1:3030"), handler.CorsMiddleware(r))
```     
    4 Save the file to take effect.

### FrontEnd
    1. Navigate to project directory and find "CompanyXYZ_Front" folder.
    2. Edit "APIURL.js".
    3. By defaut the url path is set to localhost and port 8080
```bash
      http://localhost:8080
```     
    4. Change the url if needed'
    5. Save the file to take effect.
## Dependencies

Update the Backend project Dependencies

```bash
  go mod tidy
```
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
  cd Montani
```

    2. To test, run the following command.

```bash
  go test -v ./handler
```

```bash
  go test -v ./service
```
