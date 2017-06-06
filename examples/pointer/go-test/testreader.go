package test

import "io"

type TestReader struct {
}

func (t *TestReader) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		p[i] = 'a'
	}
	return len(p), io.EOF
}
