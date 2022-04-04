FROM golang:1.18.0-bullseye
#Working directory in docker
WORKDIR $GOPATH/src/gin_docker
#Synchronize the current directory to the docker working directory, or configure only the required directories and files (configuration directory, compiled program, etc.)
ADD . ./
#For well-known reasons, some packages will have download timeouts. The proxy service of go module is also used in docker
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOPROXY="https://goproxy.io"
#Specify the file name after compilation. You can use the default file name without setting. The last step is to execute the file name
RUN go build -o gin_docker .
EXPOSE 8080
#This is consistent with the compiled file name
ENTRYPOINT  ["./gin_docker"]