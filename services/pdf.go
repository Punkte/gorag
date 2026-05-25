package services

import (
	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

func ExtractText(p string) (string, error) {
	f, r, err := pdf.Open(p)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func ChunkText(text string, size int, overlap int) []string {
	var tab []string
	for i := 0; i < len(text); i += size - overlap {
		tab = append(tab, text[i:min(size+i, len(text))])
	}
	return tab
}
