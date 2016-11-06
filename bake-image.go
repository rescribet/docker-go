package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func defaultScratchfile() string {
	return os.ExpandEnv("$GOPATH/src/github.com/fletcher91/docker-go/Dockerfile.scratch")
}

func transformScratchDockerfile(outFile string) {
	if defaultScratchfile() == *scratchfile {
		dockerfile, err := ioutil.ReadFile(*scratchfile)
		handleErr(err)
		replaced := bytes.Replace(
			dockerfile,
			[]byte("${EXEC_NAME}"),
			[]byte(*execName),
			2,
		)
		replaced = bytes.Replace(
			replaced,
			[]byte("${ARTIFACTS_PATH}"),
			[]byte(*artifacts),
			2,
		)
		ioutil.WriteFile(outFile, replaced, 0777)
	}
}

func bakeImage() {
	artPath := fmt.Sprintf("%s%s/%s", os.Getenv("GOPATH"), *projectName, *artifacts)
	tScratchFile := fmt.Sprintf("%s/Dockerfile.scratch.final", artPath)
	transformScratchDockerfile(tScratchFile)
	bakeCmd := exec.Command(
		"docker", "build",
		"-t", *imageName,
		"-f", tScratchFile,
		fmt.Sprintf("%s%s", os.Getenv("GOPATH"), *projectName),
	)
	executeCommand(bakeCmd)
	fmt.Println(*imageName)
}
