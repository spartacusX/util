package archive

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type Gzip struct {
	Src string
	Dst string
}

func (g *Gzip) Compress() error {
	logHeader := "F: Gzip.Compress, A: , M:"
	log.Infof("%s Compressing src=%s to dst=%s", logHeader, g.Src, g.Dst)

	rf, err := os.Open(g.Src)
	if err != nil {
		log.Errorf("%s err: %s", logHeader, err.Error())
		return err
	}

	defer rf.Close()

	fbuf, err := os.Create(g.Dst)
	if err != nil {
		log.Errorf("%s err: %s", logHeader, err.Error())
		return err
	}
	defer fbuf.Close()

	if data, err := ioutil.ReadAll(rf); err == nil {
		gw := gzip.NewWriter(fbuf)
		defer gw.Close()
		if _, err = gw.Write(data); err != nil {
			log.Errorf("%s writer data to %s failed, err: %s", logHeader, g.Dst, err.Error())
			return err
		}
	} else {
		log.Errorf("%s read data from %s failed, err: %s", logHeader, g.Src, err.Error())
		return err
	}

	return nil
}

func (g *Gzip) Uncompress() error {
	logHeader := "F: Gzip.Uncompress, A: , M:"
	log.Infof("%s Uncompressing src=%s to dst=%s", logHeader, g.Src, g.Dst)

	gzf, err := os.Open(g.Src)
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

	fw, err := os.Create(g.Dst)
	if err != nil {
		log.Errorf("%s err: %s", logHeader, err.Error())
		return err
	}
	defer fw.Close()

	if _, err = io.Copy(fw, gzr); err != nil {
		log.Errorf("%s copy data to %s failed.", logHeader, g.Dst)
		return err
	}

	return nil
}
