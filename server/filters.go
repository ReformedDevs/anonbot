package server

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/flosch/pongo2"
)

func filterMD5(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	h := md5.New()
	io.WriteString(h, in.String())
	return pongo2.AsValue(fmt.Sprintf("%x", h.Sum(nil))), nil
}

func init() {
	pongo2.RegisterFilter("md5", filterMD5)
}
