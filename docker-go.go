package main

import (
	"bytes"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	artifacts   = flag.String("artifacts", "go-artifacts", "The name of the artifacts folder")
	buildfile   = flag.String("build-file", defaultBuildfile(), "Location of the dockerfile to use for building the executable")
	scratchfile = flag.String("scratch-file", defaultScratchfile(), "Location of the dockerfile to use for building the final image")
	imageName   = flag.String("name", filepath.Base(*projectName), "The name of the image to bake")
	projectName = flag.String("project", curProject(), "The (import) name of the project to be built")
	execName    = flag.String("filename", "main", "The name of the executable to be generated")
	verbose     = flag.Bool("v", false, "Verbose output")
)

func curProject() string {
	wDir, err := os.Getwd()
	handleErr(err)
	return strings.SplitAfter(wDir, os.Getenv("GOPATH"))[1]
}

func curDirName() string {
	return os.ExpandEnv(fmt.Sprintf("${GOPATH}/%s", curProject()))
}

func executeCommand(cmd *exec.Cmd) {
	log.Debug("%+v\n", cmd)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if stderr.Len() > 0 {
		log.Error(fmt.Sprintf("%q\n", stderr.String()))
	}
	handleErr(err)
	log.Debug(fmt.Sprintf("%q\n", stdout.String()))
}

func fullProjectPath() string {
	return strings.SplitAfter(*projectName, "src/")[1]
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func cleanup() {
	tempPath := filepath.Join(curDirName(), *artifacts)
	log.Debug(fmt.Sprintf("Removing: %s", tempPath))
	os.RemoveAll(tempPath)
	cmd := exec.Command(
		"docker", "rmi", "-f",
		buildImageName(),
	)
	cmd.Run()
	log.Debug("Finished cleanup")
}

func main() {
	defer cleanup()
	flag.Parse()
	if *verbose == true {
		log.SetLevel(log.DebugLevel)
		log.Info("Showing debugging info")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	*projectName = os.ExpandEnv(*projectName)

	createBuildImage()
	generateArtifacts()
	bakeImage()
}
