package main

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
)

func buildImageName() string {
	return fmt.Sprintf("build-img-%s", *imageName)
}

func buildfileName() string {
	fpath := *buildfile
	if defaultBuildfile() == *buildfile {
		fpath = fmt.Sprintf("%s.final", *buildfile)
	}
	return os.ExpandEnv(fpath)
}

func defaultBuildfile() string {
	return os.ExpandEnv("$GOPATH/src/github.com/fletcher91/docker-go/Dockerfile.build")
}

func transformBuildDockerfile(outFile string) {
	if defaultBuildfile() == *buildfile {
		dockerfile, err := ioutil.ReadFile(*buildfile)
		handleErr(err)
		replaced := bytes.Replace(
			dockerfile,
			[]byte("${PROJECT_DIR}"),
			[]byte(fmt.Sprintf("/go/%s", curProject())),
			1,
		)
		log.Debug("Buildfile: %s\n", outFile)
		ioutil.WriteFile(outFile, replaced, 0777)
	}
}

func createBuildImage() {
	err := os.Mkdir(*artifacts, 0770)
	handleErr(err)
	artPath := fmt.Sprintf("%s%s/%s", os.Getenv("GOPATH"), *projectName, *artifacts)
	tBuildFile := fmt.Sprintf("%s/Dockerfile.build.final", artPath)
	transformBuildDockerfile(tBuildFile)
	if _, err = os.Stat(tBuildFile); os.IsNotExist(err) {
		handleErr(err)
	}
	buildCmd := exec.Command(
		"docker", "build",
		"-t", buildImageName(),
		"-f", tBuildFile,
		os.Getenv("GOPATH"),
	)
	executeCommand(buildCmd)
}
