package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// import (
// 	"archive/tar"
// 	"compress/gzip"
// 	"fmt"
// 	"io"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	log "github.com/sirupsen/logrus"
// )

// // Targz takes a source and variable writers and walks 'source' writing each file
// // found to the tar writer; the purpose for accepting multiple writers is to allow
// // for multiple outputs (for example a file, or md5 hash)
// func Compress(src string, filter []string, writers ...io.Writer) error {
// 	logHeader := fmt.Sprintf("F: Targz, A: src=%s, filter=%v, M:", src, filter)
// 	log.Infof("%s Called", logHeader)

// 	if _, err := os.Stat(src); err != nil {
// 		log.Errorf("%s error: %s", logHeader, err.Error())
// 		return err
// 	}

// 	mw := io.MultiWriter(writers...)

// 	gzw := gzip.NewWriter(mw)
// 	defer gzw.Close()

// 	tw := tar.NewWriter(gzw)
// 	defer tw.Close()

// 	// walk path
// 	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
// 		if info.Mode().IsDir() {
// 			return nil
// 		}

// 		log.Infof("%s path=%s", logHeader, path)
// 		validPath := false
// 		for _, subPath := range filter {
// 			if strings.HasPrefix(path, src+"/"+subPath) {
// 				validPath = true
// 				break
// 			}
// 		}

// 		if !validPath {
// 			return nil
// 		}

// 		//log.Infof("%s compressing file: %s", logHeader, path)
// 		fr, err := os.Open(path)
// 		if err != nil {
// 			log.Errorf("%s error: %s", logHeader, err.Error())
// 			return err
// 		}
// 		defer fr.Close()

// 		h, err := tar.FileInfoHeader(info, path)
// 		if err != nil {
// 			log.Errorf("%s can NOT get tar file header. error: %s", logHeader, err.Error())
// 			return err
// 		}

// 		h.Name = filepath.Base(src) + path[len(src):]
// 		if err = tw.WriteHeader(h); err != nil {
// 			log.Errorf("%s can NOT write tar file header. error: %s", logHeader, err.Error())
// 			return err
// 		}

// 		if _, err := io.Copy(tw, fr); err != nil {
// 			log.Errorf("%s can NOT write to tar file. error: %s", logHeader, err.Error())
// 			return err
// 		}

// 		return nil
// 	})

// }

func Uncompress(targetDir string, tarRdr *tar.Reader) error {
	logHeader := fmt.Sprintf("F: SMAunTar, A: targetDir=%s, M:", targetDir)
	log.Infof("%s Called", logHeader)

	_, err := os.Stat(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(targetDir, 0755); err != nil {
				log.Errorf("%s Failed to create targetDir: %s, Error: %s", logHeader, targetDir, err.Error())
				return err
			}
		}
	}

	for {
		header, err := tarRdr.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("%s Error: %s", logHeader, err.Error())
			return err
		}

		target := filepath.Join(strings.TrimRight(targetDir, "/"), "/", header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(target, header.FileInfo().Mode())
			if err != nil {
				log.Errorf("%s Error: %s", logHeader, err.Error())
				return err
			}

		case tar.TypeReg:
			// make sure we have the directory created
			//err := os.MkdirAll(filepath.Dir(target), header.FileInfo().Mode())
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
