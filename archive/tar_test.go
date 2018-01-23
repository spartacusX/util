package archive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArchive(t *testing.T) {
	tar := Tar{
		Name:    "testdata.tar",
		SrcPath: "./testdata/testarchive",
		DstPath: "./testdata/",
	}

	err := tar.Archive()
	if err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	if _, err = os.Stat(tar.DstPath); err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	} else {
		os.Remove(tar.DstPath)
	}
}

func TestUnArchive(t *testing.T) {
	tar := Tar{
		Name:    "testdata.tar",
		SrcPath: "./testdata/testunarchive/",
		DstPath: "./testdata/testunarchive/",
	}

	err := tar.UnArchive()
	if err != nil {
		t.Errorf("expected: no err, actual: err=%s", err.Error())
	}

	defer os.RemoveAll(tar.DstPath + "testarchive")

	expectedPath := []string{
		"testunarchive",
		"testdata.tar",
		"testarchive",
		"testarchive/data1",
		"testarchive/data2",
		"testarchive/folder1",
		"testarchive/folder1/f1d1",
		"testarchive/folder1/f1d2",
		"testarchive/folder2",
		"testarchive/folder2/f1d1",
		"testarchive/folder2/f1d2",
	}

	filepath.Walk(tar.DstPath, func(actualPath string, info os.FileInfo, err error) error {
		failed := true
		for _, path := range expectedPath {
			if strings.Contains(actualPath, path) {
				failed = false
				break
			}
		}
		if failed {
			t.Errorf("found invalid path: %s", actualPath)
		}
		return nil
	})
}
