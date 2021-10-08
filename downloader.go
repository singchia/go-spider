package spider

import (
	"cpf_server/library/cpf_common/common"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Downloader interface {
	//url, header, reader, suffix
	//url location, header location, body location
	Download(*url.URL, http.Header, io.Reader, string) (*string, *string, *string, error)
}

type FileDownloader struct {
	path string
}

func NewFileDownloader(path string) *FileDownloader {
	return &FileDownloader{path: path}
}

func (fd *FileDownloader) Download(u *url.URL, header http.Header, reader io.Reader, suffix string) (*string, *string, *string, error) {

	m := md5.New()
	m.Write([]byte(u.String()))
	prefix := hex.EncodeToString(m.Sum(nil))
	bodyFile := fmt.Sprintf("%s%s", prefix, suffix)

	bodyPath := filepath.Join(fd.path, bodyFile)
	if err := writeFromReader(reader, bodyPath); err != nil {
		return nil, nil, nil, err
	}

	headerPath := fmt.Sprintf("%s.%s", bodyPath, "hdr")
	if err := writeHttpHeader(header, headerPath); err != nil {
		return nil, nil, nil, err
	}

	urlPath := fmt.Sprintf("%s.%s", bodyPath, "url")
	if err := writeFromBytes([]byte(u.String()), urlPath); err != nil {
		return nil, nil, nil, err
	}

	return &urlPath, &headerPath, &bodyPath, nil
}

func writeHttpHeader(header http.Header, path string) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	return header.Write(fd)
}

func writeFromReader(reader io.Reader, path string) error {

	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = common.ReadAndWrite(reader, fd)
	if err != nil {
		return err
	}
	return nil
}

func writeFromBytes(bytes []byte, path string) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = common.Write(bytes, fd)
	if err != nil {
		return err
	}
	return nil
}
