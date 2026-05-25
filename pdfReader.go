package main

import (
	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

func extractText(p string) (string, error) {
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
