# The go gin task APIs on docker 

This project doing the implementation of CRUD APIs by gin

## Getting Started

### Summary
1. Get the latest source code from the [GitHub page](https://github.com/BinHsu/go_gin_docker).
2. Generate and Run docker
3. Run the test

### Step 1: Get the latest source code

Make sure you can run `git` with an Internet connection.

```shell
$ git clone https://github.com/BinHsu/go_gin_docker.git && cd go_gin_docker
```

### Step 2: Build docker

To generate and run the docker image

```shell
# build image
$ docker build -t gin_docker .
# run docker
$ docker run --name gin_docker -p 8080:8080 -d gin_docker
```

### Step 3: Run the tests

Now we could run test with following commands:

```bash
$ go test -v
```

## Note

The default port in this project is 8080
