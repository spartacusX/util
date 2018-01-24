package archive

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type Case struct {
	Name string
	Src  string
	err  error
}

func TestTarUnTar(t *testing.T) {
	cases := []Case{
		{"case1", "./testdata/tar", nil},
		{"case2", "./testdata/tar", nil},
		{"case3", "./testdata/tar", nil},
	}

	for _, c := range cases {
		srcDir := filepath.Join(c.Src, c.Name)
		tempDir := filepath.Join(c.Src, "temp")
		tarPath := srcDir + ".tar"

		if err := Tar(srcDir); err != c.err {
			t.Errorf("expected: %v, actual: err=%s", c.err, err.Error())
		} else {
			if _, err := os.Stat(tarPath); err != nil {
				t.Errorf("expected: no err, actual: err=%s", err.Error())
			}
		}

		if err := os.Mkdir(tempDir, 0755); err != nil {
			t.Errorf("Failed to create directory: %s. err: %s", tempDir, err.Error())
		}

		if err := UnTar(tarPath, tempDir); err != nil {
			t.Errorf("Failed to UnTar, err: %s", err.Error())
		}

		os.Remove(tarPath)

		if CompareFile(srcDir, filepath.Join(c.Src, "temp", c.Name)) == false {
			t.Errorf("Difference found for folder: %s between before and after archive", srcDir)
		}

		os.RemoveAll(tempDir)
	}
}

func CompareFile(src, dst string) bool {
	sf, err := os.Stat(src)
	if err != nil {
		return false
	}
	df, err := os.Stat(dst)
	if err != nil {
		return false
	}

	if !sf.IsDir() && !df.IsDir() {
		return sf.Name() == df.Name() && sf.Mode() == df.Mode()
	}

	srcFiles, err := ioutil.ReadDir(src)
	if err != nil {
		return false
	}
	dstFiles, err := ioutil.ReadDir(dst)
	if err != nil {
		return false
	}

	for _, srcFile := range srcFiles {
		found := false
		for _, dstFile := range dstFiles {
			if srcFile.Name() == dstFile.Name() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
