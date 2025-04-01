package compression

import (
    "compress/gzip"
    "io"
)

type GzipCompressor struct{}

func NewGzipCompressor() *GzipCompressor {
    return &GzipCompressor{}
}

func (g *GzipCompressor) Compress(data io.Reader) (io.Reader, error) {
    pr, pw := io.Pipe()
    go func() {
        gw := gzip.NewWriter(pw)
        _, err := io.Copy(gw, data)
        gw.Close()
        pw.CloseWithError(err)
    }()
    return pr, nil
}

func (g *GzipCompressor) Decompress(data io.Reader) (io.Reader, error) {
    return gzip.NewReader(data)
}