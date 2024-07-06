package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/hop-/pdf-service/internal/generators"
)

type DocData struct {
	Data map[string]any `json:"data" validate:"required"`
}

type Doc struct {
	Content string `json:"content"`
}

func NewRouter() http.Handler {
	router := http.NewServeMux()

	router.Handle("/docs/", http.StripPrefix("/docs", newDocsRouter()))

	return router
}

func newDocsRouter() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /{type}", newDocHandler)
	router.HandleFunc("POST /{type}/pdf", serveNewDocFileHandler)

	return router
}

func newDocHandler(w http.ResponseWriter, r *http.Request) {
	docType := r.PathValue("type")
	var req DocData
	json.NewDecoder(r.Body).Decode(&req)

	content, err := generators.GetConcurrentPdfGenerator().GenerateBase64(docType, req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	doc := Doc{
		Content: content,
	}

	res, err := json.Marshal(doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(res)
}

func serveNewDocFileHandler(w http.ResponseWriter, r *http.Request) {
	docType := r.PathValue("type")
	var req DocData
	json.NewDecoder(r.Body).Decode(&req)

	content, err := generators.GetConcurrentPdfGenerator().Generate(docType, req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	http.ServeContent(w, r, docType+".pdf", time.Now(), bytes.NewReader(content))
}
