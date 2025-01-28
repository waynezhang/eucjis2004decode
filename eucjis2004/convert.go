package eucjis2004

import (
	"bytes"
)

func Convert(in []byte, out *bytes.Buffer) error {
	dec := EUCJIS2004Decoder{}

	buf := make([]byte, 4096)

	idx := 0
	for {
		nDst, nSrc, err := dec.Transform(buf, in[idx:], true)
		if err != nil {
			return err
		}
		if nSrc == 0 {
			break
		}
		out.Write(buf[:nDst])
		idx += nSrc
	}

	return nil
}
