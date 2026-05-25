package main

func chunckText(text string, size int, overlap int) (tab []string) {
	for i := 0; i < len(text); i += size - overlap {
		tab = append(tab, text[i:min(size+i, len(text))])
	}

	return tab
}
