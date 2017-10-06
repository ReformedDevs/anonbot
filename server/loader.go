package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// b0xLoader provides a Pongo2 loader for b0x template files.
type b0xLoader struct{}

// Abs returns the absolute path to a template file.
func (l b0xLoader) Abs(base, name string) string {
	return name
}

// Get retrieves a reader for the specified path.
func (l b0xLoader) Get(path string) (io.Reader, error) {
	f, err := FS.OpenFile(
		CTX,
		fmt.Sprintf("templates/%s", path),
		os.O_RDONLY,
		0,
	)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
