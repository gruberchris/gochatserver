# Go Chat Server

A better TCP socket chat server implemented using Go

## Connecting
```shell
nc localhost 5000
```

## Docker

Build the container image

~~~
docker build -t gruberchris/gochatserver:dev .
~~~

Execute the container
~~~
docker run -d -p 5000:5000 --name gochatserver gruberchris/gochatserver:dev
~~~