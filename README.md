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

The final image is built from the [scratch image](https://hub.docker.com/_/scratch/), therefore they are generally between 5-20MB depending on the size of your (statically built) executable.

This project is based on the article [Building Minimal Docker Containers for Go Applications](https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/) by Nick Gauthier.
