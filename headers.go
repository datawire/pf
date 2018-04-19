// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
)

var re = regexp.MustCompile("root:xnu-([0-9.]+)")

func GetKernelVersion() (result string, err error) {
	cmd := exec.Command("uname", "-v")
	out, err := cmd.CombinedOutput()
	if err != nil { return }
	matches := re.FindSubmatch(out)
	if len(matches) < 2 {
		err = fmt.Errorf("cannot find version in uname output: %s", out)
		return
	}
	result = string(matches[1])
	return 
}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

type Header struct {
	Base string
	Name string
	Prefix string
}

func (h *Header) Fetch(version string) (err error) {
	if exists, err := Exists(h.Name); exists || err != nil {
		return err
	}

	url := fmt.Sprintf("https://opensource.apple.com/source/xnu/xnu-%s/%s%s", version, h.Base, h.Name)
	r, err := http.Get(url)
	if err != nil { return }
	defer r.Body.Close()
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
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("This package will only work on Darwin.")
	}

	version, err := GetKernelVersion()
	if err != nil { return }

	fmt.Printf("Fetching headers for %s, %s.\n", runtime.GOOS, version)

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
