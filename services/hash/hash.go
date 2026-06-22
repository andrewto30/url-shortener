package hash

import "io"

type Generator struct {
	randSource io.Reader
}

func New(src io.Reader) *Generator {
	return &Generator{randSource: src}
}

func (g *Generator) Key() (string, error) {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	var buf [4]byte
	if _, err := g.randSource.Read(buf[:]); err != nil {
		return "", err
	}

	out := make([]byte, 4)
	for i, b := range buf {
		out[i] = alphabet[b%62]
	}

	return string(out), nil
}
