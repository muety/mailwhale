package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

type SPAHandler struct {
	StaticPath      string
	IndexPath       string
	ReplaceBasePath string
	NoCache         bool
	indexContent    []byte
}

func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.StaticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) || r.URL.Path == "/" {
		if len(h.indexContent) == 0 || h.NoCache {
			h.loadIndex()
		}

		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(h.indexContent)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}

// dirty little hack to replace base attribute in html according to instance's public url
func (h *SPAHandler) loadIndex() {
	raw, err := os.ReadFile(filepath.Join(h.StaticPath, h.IndexPath))
	if err != nil {
		panic(err)
	}
	html := string(raw)
	if h.ReplaceBasePath != "" {
		pattern := regexp.MustCompile(`<base href="(.*)"`)
		html = pattern.ReplaceAllString(html, fmt.Sprintf(`<base href="%s"`, h.ReplaceBasePath))
	}
	h.indexContent = []byte(html)
}
