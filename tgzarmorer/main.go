package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal("len(args) != 2")
	}
	path := os.Args[1]

	log.Printf("%s: MOCKING", path)

	if err := mock(path); err != nil {
		log.Fatalf("%s: FAILED - %s", path, err.Error())
	}

	log.Printf("%s: DONE", path)
}

func mock(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter, err := gzip.NewWriterLevel(os.Stdout, gzip.NoCompression)
	if err != nil {
		return err
	}
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	tarHeader := tar.Header{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
		Mode: int64(fileInfo.Mode()),
	}
	if err := tarWriter.WriteHeader(&tarHeader); err != nil {
		return err
	}

	if _, err := io.Copy(tarWriter, file); err != nil {
		return err
	}

	return nil
}
