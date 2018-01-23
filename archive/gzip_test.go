package archive

import (
	"os"
	"testing"
)

func TestCompress(t *testing.T) {
	gzip := Gzip{
		Src: "./testdata/gzip/compress/file1",
		Dst: "./testdata/gzip/compress/file1.gz",
	}

	err := gzip.Compress()
	if err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	_, err = os.Stat(gzip.Dst)
	if err != nil {
		t.Errorf("expected: %s exists, actual: err=%s", gzip.Dst, err.Error())
	}

	os.Remove(gzip.Dst)
}

func TestUnCompress(t *testing.T) {
	gzip := Gzip{
		Src: "./testdata/gzip/uncompress/file1.gz",
		Dst: "./testdata/gzip/uncompress/file1",
	}

	err := gzip.Uncompress()
	if err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	if _, err = os.Stat(gzip.Dst); err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}
}
