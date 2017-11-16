package output

import (
	"io"
)

type prependWriter struct {
	prefix []byte
	wr     io.Writer
}

func (p *prependWriter) Write(b []byte) (n int, err error) {
	toWrite := append(p.prefix, b...)
	return p.wr.Write(toWrite)
}
