package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/opencontainers/runtime-tools/validation/util"
)

func checkReadonlyPaths() error {
	g, err := util.GetDefaultGenerator()
	if err != nil {
		return err
	}

	readonlyDir := "readonly-dir"
	readonlySubDir := "readonly-subdir"
	readonlyFile := "readonly-file"

	readonlyDirTop := filepath.Join("/", readonlyDir)
	readonlyFileTop := filepath.Join("/", readonlyFile)

	readonlyDirSub := filepath.Join(readonlyDirTop, readonlySubDir)
	readonlyFileSub := filepath.Join(readonlyDirTop, readonlyFile)
	readonlyFileSubSub := filepath.Join(readonlyDirSub, readonlyFile)

	g.AddLinuxReadonlyPaths(readonlyDirTop)
	g.AddLinuxReadonlyPaths(readonlyFileTop)
	g.AddLinuxReadonlyPaths(readonlyDirSub)
	g.AddLinuxReadonlyPaths(readonlyFileSub)
	g.AddLinuxReadonlyPaths(readonlyFileSubSub)
	err = util.RuntimeInsideValidate(g, func(path string) error {
		testDir := filepath.Join(path, readonlyDirSub)
		err = os.MkdirAll(testDir, 0777)
		if err != nil {
			return err
		}
		// create a temp file to make testDir non-empty
		tmpfile, err := ioutil.TempFile(testDir, "tmp")
		if err != nil {
			return err
		}
		defer os.Remove(tmpfile.Name())

		// runtimetest cannot check the readability of empty files, so
		// write something.
		testSubSubFile := filepath.Join(path, readonlyFileSubSub)
		if err := ioutil.WriteFile(testSubSubFile, []byte("immutable"), 0777); err != nil {
			return err
		}

		testSubFile := filepath.Join(path, readonlyFileSub)
		if err := ioutil.WriteFile(testSubFile, []byte("immutable"), 0777); err != nil {
			return err
		}

		testFile := filepath.Join(path, readonlyFile)
		return ioutil.WriteFile(testFile, []byte("immutable"), 0777)
	})
	return err
}

func main() {
	if err := checkReadonlyPaths(); err != nil {
		util.Fatal(err)
	}

}
