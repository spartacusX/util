package archive

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGzipUnGzip(t *testing.T) {
	cases := []string{
		"./testdata/gzip/case1",
		"./testdata/gzip/case2.tar",
	}

	for _, c := range cases {
		err := Gzip(c)
		if err != nil {
			t.Errorf("expected: no err, actual: err=%s", err.Error())
		}

		gzPath := c + ".gz"
		_, err = os.Stat(gzPath)
		if err != nil {
			t.Errorf("expected: %s exists, actual: err=%s", gzPath, err.Error())
		}

		tempDir := filepath.Join(filepath.Dir(c), "temp")
		if err := os.Mkdir(tempDir, 0755); err != nil {
			t.Errorf("Failed to create directory: %s, err=%s", tempDir, err.Error())
		}

		err = Ungzip(gzPath, tempDir)
		if err != nil {
			t.Errorf("expected: no err, actual: err=%s", err.Error())
		}

		os.Remove(gzPath)

		if _, err = os.Stat(filepath.Join(tempDir, filepath.Base(c))); err != nil {
			t.Errorf("expected: no err, actual: err=%s", err.Error())
		}
		os.RemoveAll(tempDir)
	}
}
