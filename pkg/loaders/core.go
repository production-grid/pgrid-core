package loaders

import (
	"io"
	"io/ioutil"
)

// ResourceLoader defines the contract for loading non code resources.
type ResourceLoader interface {
	Reader(path string) (io.Reader, error)
	Bytes(path string) ([]byte, error)
	String(path string) (string, error)
}

// Bytes is a convenience method that all ResourceLoader implementations
// can use to load a resource as bytes.
func Bytes(loader ResourceLoader, path string) ([]byte, error) {

	reader, err := loader.Reader(path)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(reader)

}

// String is a convenience method that all ResourceLoader implementations
// can use to load a resource as a string.
func String(loader ResourceLoader, path string) (string, error) {

	bytes, err := Bytes(loader, path)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
