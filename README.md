# docker-go

A CLI tool to dockerize your runnable Go app from scratch.

## Usage
Install this project
```
$ go get github.com/fletcher91/docker-go
$ cd $GOPATH/github.com/fletcher91/docker-go
$ go install
```

When having downloaded an executable go project, run docker-go in the directory:
`$ docker-go`

After which you can run your app in docker:
`$ docker run -t github.com/fletcher91/docker-go:latest`
