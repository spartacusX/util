package archive

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Tar struct {
	Name    string //archive file name
	SrcPath string //directory to be archived
	DstPath string //directory to place the archive file
}

func (t *Tar) Archive() error {
	logHeader := "F: Tar.Archive, A: , M:"
	log.Infof("%s archiving Name=%s Src=%s, Dst=%s", logHeader, t.Name, t.SrcPath, t.DstPath)

	if _, err := os.Stat(t.SrcPath); err != nil {
		log.Errorf("%s error: %s", logHeader, err.Error())
		return err
	}

	f, err := os.Create(strings.TrimRight(t.DstPath, "/") + "/" + t.Name)
	if err != nil {
		log.Errorf("%s error: %s", logHeader, err.Error())
		return err
	}

	defer f.Close()

	tw := tar.NewWriter(f)
	defer tw.Close()

	// walk path
	return filepath.Walk(t.SrcPath, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir() {
			return nil
		}

		log.Debugf("%s compressing file: %s", logHeader, path)
		fr, err := os.Open(path)
		if err != nil {
			log.Errorf("%s error: %s", logHeader, err.Error())
			return err
		}
		defer fr.Close()

		h, err := tar.FileInfoHeader(info, path)
		if err != nil {
			log.Errorf("%s can NOT get tar file header. error: %s", logHeader, err.Error())
			return err
		}

		h.Name = filepath.Base(t.SrcPath) + "/" + path[len(t.SrcPath)-1:]
		if err = tw.WriteHeader(h); err != nil {
			log.Errorf("%s can NOT write tar file header. error: %s", logHeader, err.Error())
			return err
		}

		if _, err := io.Copy(tw, fr); err != nil {
			log.Errorf("%s can NOT write to tar file. error: %s", logHeader, err.Error())
			return err
		}

		return nil
	})
}

func (t *Tar) UnArchive() error {
	logHeader := "F: Tar.UnArchive, A: , M:"
	log.Infof("%s extracting files in t.SrcPath=%s Name=%s to t.DstPath=%s", logHeader, t.SrcPath, t.Name, t.DstPath)

	f, err := os.Open(strings.TrimRight(t.SrcPath, "/") + "/" + t.Name)
	if err != nil {
		log.Errorf("%s error: %s", logHeader, err.Error())
		return err
	}

	defer f.Close()

	if _, err := os.Stat(t.DstPath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(t.DstPath, 0755); err != nil {
				log.Errorf("%s Failed to create targetDir: %s, Error: %s", logHeader, t.DstPath, err.Error())
				return err
			}
		}
	}

	tarRdr := tar.NewReader(f)
	for {
		header, err := tarRdr.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("%s Error: %s", logHeader, err.Error())
			return err
		}

		target := filepath.Join(strings.TrimRight(t.DstPath, "/"), "/", header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(target, header.FileInfo().Mode())
			if err != nil {
				log.Errorf("%s Error: %s", logHeader, err.Error())
				return err
			}

		case tar.TypeReg:
			err := os.MkdirAll(filepath.Dir(target), 0755)
			if err != nil {
				log.Errorf("%s Error: %s", logHeader, err.Error())
				return err
			}

			// create the file
			file, err := os.Create(target)
			if err != nil {
				log.Errorf("%s Error: %s", logHeader, err.Error())
				return err
			}

			if _, err := io.Copy(file, tarRdr); err != nil {
				log.Errorf("%s Error: %s", logHeader, err.Error())
				return err
			}

			file.Close()
		}
	}
	return nil
}
