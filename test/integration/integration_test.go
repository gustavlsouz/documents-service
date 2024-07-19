package integration

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/gustavlsouz/documents-service/pkg"
)

func StartContainers() error {
	cmd := exec.Command("make", "init")
	var out bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuffer

	err := cmd.Run()

	if err != nil {
		fmt.Printf("output: %s\n", errBuffer.String())
		fmt.Printf("error: %v\n", err)
		return err
	}

	fmt.Printf("output: %s\n", out.String())

	return nil
}

func TestMain(m *testing.M) {
	startedSuccessfullyChan := make(chan bool)

	err := StartContainers()
	if err != nil {
		panic("containers not started")
	}

	go pkg.Start(startedSuccessfullyChan, "../../.env.local", "../../deployments/migrations")

	startedSuccessfully := <-startedSuccessfullyChan

	if !startedSuccessfully {
		panic("server not started")
	}

	m.Run()
}
