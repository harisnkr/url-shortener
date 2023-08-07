# url-shortener

url-shortener serves simple URL shortener written in Go. This project aims to provide a convenient way to shorten long
URLs, making them more manageable and shareable.

## Features

- Shorten long URLs into easy-to-remember short links using UUIDv4.
- Gin HTTP framework for quick performance
- Dockerised Golang server and Redis server

## Usage

Pre-requisites: **Golang 1.20, Docker Desktop + docker-compose CLI**

1. Clone the repository in your local machine:
    ```bash
   git clone https://github.com/harisnkr/url-shortener && cd url-shortener
   ```


2. Run Redis and url-shortener service:
   ```shell
   docker-compose up
   ```