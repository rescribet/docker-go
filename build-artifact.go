package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"os/exec"
	"time"
)

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func buildArtifact(artImageName string, fileName string, goos string) {
	log.Debug(fmt.Sprintf("PROJECT_DIR=/go%s", *projectName))
	buildCmd := exec.Command(
		"docker", "run",
		"--name", artImageName,
		"-e", "CGO_ENABLED=0",
		"-e", fmt.Sprintf("GOOS=%s", goos),
		"-e", "GOARCH=amd64",
		buildImageName(),
		"go", "build", "-a", "-installsuffix", "cgo",
		"-o", fmt.Sprintf("/%s/%s", *artifacts, fileName),
		".",
	)
	executeCommand(buildCmd)
}

func copyArtifacts(artImageName string, fileName string) {
	cpCmd := exec.Command(
		"docker", "cp",
		fmt.Sprintf("%s:%s/%s", artImageName, *artifacts, fileName),
		fmt.Sprintf("./%s/%s", *artifacts, fileName),
	)
	executeCommand(cpCmd)
	cpCmd = exec.Command(
		"docker", "cp",
		fmt.Sprintf("%s:/etc/ssl/certs/ca-certificates.crt", artImageName),
		fmt.Sprintf("./%s/ca-certificates.crt", *artifacts),
	)
	executeCommand(cpCmd)
}

func removeArtifactContainer(artContName string) {
	log.Debug(fmt.Sprintf("Removing artifact container '%s'", artContName))
	rmCmd := exec.Command("docker", "rm", "-f", artContName)
	executeCommand(rmCmd)
}

func generateArtifact(filename string, goos string) {
	rand := randomString(10)
	artContName := fmt.Sprintf("build-img-%s-cont-%s-%s", *imageName, filename, rand)
	defer removeArtifactContainer(artContName)
	log.Debug(fmt.Sprintf("Building artifact '%s/%s'", artContName, filename))
	buildArtifact(artContName, filename, goos)
	log.Debug(fmt.Sprintf("Copying artifact '%s'", filename))
	copyArtifacts(artContName, filename)
}

func generateArtifacts() {
	generateArtifact(*execName, "linux")
}
