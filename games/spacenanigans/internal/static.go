package internal

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type staticFile struct {
	Name    string
	ModTime time.Time
	Buf     *bytes.Reader
}

func loadStatic(root string) (map[string]staticFile, error) {
	root = filepath.Clean(root)
	m := make(map[string]staticFile)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		//path = strings.TrimPrefix(path, root)
		m[path] = staticFile{
			Name:    info.Name(),
			ModTime: info.ModTime(),
			Buf:     bytes.NewReader(b),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Server) serveStatic(w http.ResponseWriter, r *http.Request, path string) {
	path = strings.TrimPrefix(path, "/")
	f, found := s.static[path]
	if !found {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	// Autodetects content-type of the file
	http.ServeContent(w, r, f.Name, f.ModTime, f.Buf)
}

func (s *Server) StaticHandler(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.serveStatic(w, r, path)
		//path = strings.TrimPrefix(path, "/")
		//f, found := s.static[path]
		//if !found {
		//http.Error(w, "not found", http.StatusNotFound)
		//return
		//}
		//// Autodetects content-type of the file
		//http.ServeContent(w, r, f.Name, f.ModTime, f.Buf)
	})
}
