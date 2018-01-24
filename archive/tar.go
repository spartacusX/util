package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Tar archive the given file/directory to be a tar in the same path
func Tar(src string) error {
	logHeader := fmt.Sprintf("F: Tar, A: src=%s, M:", src)
	log.Infof("%s Called", logHeader)

	if _, err := os.Stat(src); err != nil {
		log.Errorf("%s error: %s", logHeader, err.Error())
		return err
	}

	target := src + ".tar"
	f, err := os.Create(target)
	if err != nil {
		log.Errorf("%s Can NOT create tar %s. error: %s", logHeader, target, err.Error())
		return err
	}

	defer f.Close()

	tw := tar.NewWriter(f)
	defer tw.Close()

	// walk path
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir() {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}
			if len(files) > 0 {
				return nil
			}
		}

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

		h.Name = strings.Replace(path, strings.TrimRight(src, "/"), filepath.Base(src), 1)
		if err = tw.WriteHeader(h); err != nil {
			log.Errorf("%s can NOT write tar file header. error: %s", logHeader, err.Error())
			return err
		}

		if !info.Mode().IsDir() {
			if _, err := io.Copy(tw, fr); err != nil {
				log.Errorf("%s can NOT write to tar file. error: %s", logHeader, err.Error())
				return err
			}
		}

		return nil
	})
}

// UnTar unarchive the given tar file to path dst
func UnTar(src string, dst string) error {
	logHeader := fmt.Sprintf("F: UnTar, A: src=%s, M:", src)
	log.Infof("%s Called.", logHeader)

	f, err := os.Open(src)
	if err != nil {
		log.Errorf("%s Can NOT open tar: %s. error: %s", logHeader, src, err.Error())
		return err
	}

	defer f.Close()

	if _, err := os.Stat(dst); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(dst, 0755); err != nil {
				log.Errorf("%s Failed to create dst: %s, Error: %s", logHeader, dst, err.Error())
				return err
			}
		} else {
			log.Errorf("%s Got error when stat: %s. error: %s", logHeader, dst, err.Error())
			return err
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

		target := filepath.Join(strings.TrimRight(dst, "/"), "/", header.Name)

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
