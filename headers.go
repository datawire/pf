// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Header struct {
	Base string
	Name string
	Prefix string
}

func (h *Header) Fetch(version string) (err error) {
	url := fmt.Sprintf("https://opensource.apple.com/source/xnu/xnu-%s/%s%s", version, h.Base, h.Name)
	r, err := http.Get(url)
	if err != nil { return }
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return fmt.Errorf("error fetching %s for %s: %s", h.Name, version, r.Status)
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil { return }

	err = os.MkdirAll(filepath.Dir(h.Name), os.ModePerm);
	if err != nil { return }

	f, err := os.Create(h.Name)
	if err != nil { return }

	_, err = f.Write([]byte(h.Prefix))
	if err != nil { return }
	_, err = f.Write(bytes)

	return
}

func doit() (err error) {
	// For now we just vendor in the headers. If apple makes a
	// binary incompatible change at some point in the future we
	// may need a more dynamic strategy.
	version := "3789.70.16"

	fmt.Printf("Fetching headers for xnu-%s.\n", version)

	for _, hdr := range []Header {
		{"bsd/", "net/pfvar.h", "#define PRIVATE\n\n"},
		{"bsd/", "net/radix.h", ""},
		{"libkern/", "libkern/tree.h", ""},
	} {
		err = hdr.Fetch(version)
		if err != nil { return }
	}

	return
}

func main() {
	err := doit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
