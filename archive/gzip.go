package archive

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Gzip compress the given file src to src.gz in the same directory
func Gzip(src string) error {
	logHeader := fmt.Sprintf("F: Gzip, A: src=%s, M:", src)
	log.Infof("%s Called.", logHeader)

	rf, err := os.Open(src)
	if err != nil {
		log.Errorf("%s err: %s", logHeader, err.Error())
		return err
	}

	defer rf.Close()

	dst := src + ".gz"
	fbuf, err := os.Create(dst)
	if err != nil {
		log.Errorf("%s err: %s", logHeader, err.Error())
		return err
	}
	defer fbuf.Close()

	if data, err := ioutil.ReadAll(rf); err == nil {
		gw := gzip.NewWriter(fbuf)
		defer gw.Close()
		if _, err = gw.Write(data); err != nil {
			log.Errorf("%s writer data to %s failed, err: %s", logHeader, dst, err.Error())
			return err
		}
	} else {
		log.Errorf("%s read data from %s failed, err: %s", logHeader, src, err.Error())
		return err
	}

	return nil
}

// Ungzip uncompress the given src.gz to directory dst
func Ungzip(src, dst string) error {
	logHeader := fmt.Sprintf("F: Ungzip, A: src=%s dst=%s, M:", src, dst)
	log.Infof("%s Called.", logHeader)

	gzf, err := os.Open(src)
	if err != nil {
		log.Errorf("%s err: %s", logHeader, err.Error())
		return err
	}

	defer gzf.Close()

	gzr, err := gzip.NewReader(gzf)
	if err != nil {
		log.Errorf("%s err: %s", logHeader, err.Error())
		return err
	}
	defer gzr.Close()

	targetFile := filepath.Join(dst, strings.TrimRight(filepath.Base(src), ".gz"))
	fw, err := os.Create(targetFile)
	if err != nil {
		log.Errorf("%s err: %s", logHeader, err.Error())
		return err
	}
	defer fw.Close()

	if _, err = io.Copy(fw, gzr); err != nil {
		log.Errorf("%s copy data to %s failed.", logHeader, targetFile)
		return err
	}

	return nil
}
