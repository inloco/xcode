package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal("len(args) != 2")
	}
	path := os.Args[1]

	log.Printf("%s: VERIFYING", path)

	if err := verify(path); err != nil {
		log.Fatalf("%s: FAILED - %s", path, err.Error())
	}

	log.Printf("%s: PASSED", path)
}

func verify(path string) error {
	xar, err := NewXAR(path, XOFRead)
	if err != nil {
		return err
	}
	defer xar.Close()

	if err := xar.RegisterErrHandler(); err != nil {
		return err
	}

	if err := xar.Verify(); err != nil {
		return err
	}

	return nil
}
