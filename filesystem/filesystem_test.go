package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir(t *testing.T) {
	if err := CopyDir("./testdata/src_dir", "./testdata/dst_dir"); err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	srcfiles, err := ioutil.ReadDir("./testdata/src_dir")
	if err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	dstfiles, err := ioutil.ReadDir("./testdata/dst_dir")
	if err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	if len(srcfiles) != len(dstfiles) {
		t.Errorf("expected: number of files in dst_dir = %d, actual: %d", len(srcfiles), len(dstfiles))
	}

	//cleanup
	os.RemoveAll("./testdata/dst_dir")
}

func TestCopySubFolder(t *testing.T) {
	if err := CopyDir(filepath.Join("./testdata/src_dir", "dir1"), filepath.Join("./testdata/dst_dir", "dir1")); err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	srcfiles, err := ioutil.ReadDir("./testdata/src_dir/dir1")
	if err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	dstfiles, err := ioutil.ReadDir("./testdata/dst_dir/dir1")
	if err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	if len(srcfiles) != len(dstfiles) {
		t.Errorf("expected: number of files in dst_dir = %d, actual: %d", len(srcfiles), len(dstfiles))
	}

	//cleanup
	os.RemoveAll("./testdata/dst_dir")
}
