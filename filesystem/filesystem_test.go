package filesystem

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCopyDirAll(t *testing.T) {
	if err := CopyDir("./testdata/src_dir", "./testdata/dst_dir", nil); err != nil {
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

func TestCopyDirFilter(t *testing.T) {
	filters := []string{
		"testdata/src_dir/file2",
		"testdata/src_dir/dir1/d1f2",
	}

	if err := CopyDir("./testdata/src_dir", "./testdata/dst_dir", filters); err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	dstfiles, err := ioutil.ReadDir("./testdata/dst_dir")
	if err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	if 2 != len(dstfiles) {
		t.Errorf("expected: number of files in dst_dir = 2, actual: %d", len(dstfiles))
	}

	//cleanup
	os.RemoveAll("./testdata/dst_dir")
}
