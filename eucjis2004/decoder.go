package eucjis2004

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

type EUCJIS2004Decoder struct {
	transform.NopResetter
}

func (EUCJIS2004Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	r, size := rune(0), 0
loop:
	for ; nSrc < len(src); nSrc += size {
		switch c0 := src[nSrc]; {
		case c0 < 0x7F:
			r, size = rune(c0), 1

		case c0 == 0x8e:
			if nSrc+1 >= len(src) {
				if !atEOF {
					err = transform.ErrShortSrc
					break loop
				}
				r, size = utf8.RuneError, 1
				break
			}
			c1 := src[nSrc+1]
			switch {
			case c1 < 0xa1:
				r, size = utf8.RuneError, 1
			case c1 > 0xdf:
				r, size = utf8.RuneError, 2
				if c1 == 0xff {
					size = 1
				}
			default:
				r, size = rune(c1)+(0xff61-0xa1), 2
			}
		case c0 == 0x8f:
			if nSrc+2 >= len(src) {
				if !atEOF {
					err = transform.ErrShortSrc
					break loop
				}
				r, size = utf8.RuneError, 1
				if p := nSrc + 1; p < len(src) && 0xa1 <= src[p] && src[p] <= 0xfe {
					size = 2
				}
				break
			}
			c1 := src[nSrc+1]
			if c1 < 0xa1 || c1 > 0xfe {
				r, size = utf8.RuneError, 1
				break
			}
			c2 := src[nSrc+2]
			if c2 < 0xa1 || (c1 == 0xfe && c2 > 0xf6) {
				r, size = utf8.RuneError, 2
				break
			}
			r, size = utf8.RuneError, 3

			if i := int(c1)<<8 + int(c2) - 0xa1a1; i < len(eucJis2ndMap) {
				r = eucJis2ndMap[i]
			}

		case 0xa1 <= c0 && c0 <= 0xfe:
			if nSrc+1 >= len(src) {
				if !atEOF {
					err = transform.ErrShortSrc
					break loop
				}
				r, size = utf8.RuneError, 1
				break
			}
			c1 := src[nSrc+1]
			if c1 < 0xa1 || 0xfe < c1 {
				r, size = utf8.RuneError, 1
				break
			}
			r, size = utf8.RuneError, 2

			i := int(c0)<<8 + int(c1)
			if rs, ok := eucJis1stCombined[i]; ok {
				if nDst+utf8.RuneLen(rs[0]) > len(dst) {
					err = transform.ErrShortDst
					break loop
				}
				nDst += utf8.EncodeRune(dst[nDst:], rs[0])
				r, size = rs[1], 2
			} else if i-0xa1a1 < len(eucJis1stMap) {
				r = eucJis1stMap[i-0xa1a1]
			} else {
				r = utf8.RuneError
			}
		default:
			r, size = utf8.RuneError, 1
		}

		if r == utf8.RuneError {
			err = errors.New("invalid rune")
			break loop
		}

		if nDst+utf8.RuneLen(r) > len(dst) {
			err = transform.ErrShortDst
			break loop
		}
		nDst += utf8.EncodeRune(dst[nDst:], r)
	}

	return nDst, nSrc, err
}
