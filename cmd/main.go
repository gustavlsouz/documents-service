package main

import (
	"github.com/gustavlsouz/documents-service/pkg"
)

func main() {
	pkg.Start(make(chan<- bool), "../.env.local", "../deployments/migrations")
}
