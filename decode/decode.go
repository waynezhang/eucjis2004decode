package decode

import (
	"bytes"
	"errors"
	"fmt"

	tbl "github.com/waynezhang/eucjis2004decode/table"
)

func Convert(in []byte, out *bytes.Buffer) error {
	if len(in) == 0 {
		return nil
	}

	idx := 0
	max_len := len(in)

	for idx < max_len {
		b := in[idx]
		switch {
		case 0x00 <= b && b < 0x7F:
			out.WriteRune(rune(b))

			idx++
			break

		case b == 0x8E:
			if idx+1 >= max_len {
				return errors.New("Invalid code")
			}

			b2 := in[idx+1]
			if b2 >= 0xE0 {
				return errors.New("Invalid code")
			}
			r := tbl.EUC_JIS_HANKAKU_MAP[b2-0xA1]
			out.WriteRune(r)

			idx += 2
			break

		case b == 0x8F:
			if idx+2 >= max_len {
				return errors.New("Invalid code")
			}

			b2 := in[idx+1]
			if b2 < 0xA1 {
				return errors.New("Invalid code")
			}

			b3 := in[idx+2]
			if b3 < 0xA1 || (b2 == 0xFE && b3 > 0xF6) {
				return errors.New("Invalid code")
			}

			key := int32(b2)<<8 + int32(b3)
			r := tbl.EUC_JIS_2ND_MAP[key-0xA1A1]
			out.WriteRune(r)

			idx += 3
			break

		case 0xA1 <= b && b <= 0xFE:
			if idx+1 >= max_len {
				return errors.New("Invalid code")
			}

			b2 := in[idx+1]
			if b2 < 0xA1 {
				return errors.New("Invalid code")
			}

			key := int32(b)<<8 + int32(b2)
			if runes, existed := tbl.EUC_JIS_1ST_COMBINED[key]; existed {
				out.WriteRune(runes[0])
				out.WriteRune(runes[1])
			} else {
				r := tbl.EUC_JIS_1ST_MAP[key-0xA1A1]
				out.WriteRune(r)
			}

			idx += 2
			break

		default:
			return errors.New(fmt.Sprintf("Invalid code %x", b))
		}
	}

	return nil
}
