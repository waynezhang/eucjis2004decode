package eucjis2004

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	f, err := os.Open("../testdata/euc-jis-2004-with-char.txt")
	assert.Nil(t, err)
	defer f.Close()

	expected_f, err := os.Open("../testdata/expected.txt")
	assert.Nil(t, err)
	defer expected_f.Close()

	buf := bytes.NewBuffer(nil)

	src_s := bufio.NewScanner(f)
	expected_s := bufio.NewScanner(expected_f)
	idx := 0
	for expected_s.Scan() {
		buf.Reset()

		idx++
		assert.True(t, src_s.Scan())

		if idx == 216 {
			hexDump := hex.EncodeToString(src_s.Bytes())
			fmt.Println(src_s.Text(), hexDump)

			hexDump = hex.EncodeToString(expected_s.Bytes())
			fmt.Println(expected_s.Text(), hexDump)
		}

		err := Convert(src_s.Bytes(), buf)
		assert.Nil(t, err)
		assert.Equal(t, expected_s.Bytes(), buf.Bytes(), "line "+strconv.Itoa(idx))
	}
}

func TestConvert2(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	err := Convert([]byte{}, buf)
	assert.Nil(t, err)

	cases := [][]interface{}{
		{string(rune(0x0000)), []byte{0x00}},
		{string(rune(0x0001)), []byte{0x01}},
		{string(rune(0x007E)), []byte{0x7E}},

		{string(rune(0x3000)), []byte{0xA1, 0xA1}},
		{string(rune(0x32BF)), []byte{0xA8, 0xDE}},
		{string(rune(0x0000)), []byte{0xA8, 0xDF}},
		{string(rune(0x25D0)), []byte{0xA8, 0xE7}},
		{string(rune(0xFF9F)), []byte{0x8E, 0xDF}},

		{string(rune(0x20089)), []byte{0x8F, 0xA1, 0xA1}},
		{string(rune(0x2A6B2)), []byte{0x8F, 0xFE, 0xF6}},
	}

	for idx, c := range cases {
		msg := "case " + strconv.Itoa(idx)

		buf.Reset()
		err = Convert(c[1].([]byte), buf)
		assert.Nil(t, err, msg)
		assert.Equal(t, c[0].(string), buf.String(), msg)
	}
}

func TestDecode2(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	err := Convert([]byte{}, buf)
	assert.Nil(t, err)

	cases := [][]byte{
		{0x7F},
		{0xFF},

		{0xA1},
		{0xA1, 0xA0},
		{0x8E},
		{0x8E, 0xE0},

		{0x8F},
		{0x8F, 0xA0},
		{0x8F, 0xA1},
		{0x8F, 0xA0, 0xA1},
		{0x8F, 0xA1, 0xA0},
		{0x8F, 0xA0, 0xA0},

		{0x8F, 0xFE, 0xF7},
		{0x8F, 0xFE, 0xF8},
		{0x8F, 0xFF, 0x00},
	}

	for idx, c := range cases {
		msg := "case " + strconv.Itoa(idx) + " " + hex.Dump(c)

		buf.Reset()
		err = Convert(c, buf)
		assert.NotNil(t, err, msg)
	}
}
