package eucjis2004

import (
	"bufio"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/transform"
)

func TestDecode(t *testing.T) {
	f, err := os.Open("../testdata/euc-jis-2004-with-char.txt")
	assert.Nil(t, err)
	defer f.Close()

	expected_f, err := os.Open("../testdata/expected.txt")
	assert.Nil(t, err)
	defer expected_f.Close()

	decoder := &EUCJIS2004Decoder{}

	src_s := bufio.NewScanner(transform.NewReader(f, decoder))
	expected_s := bufio.NewScanner(expected_f)
	idx := 0
	for expected_s.Scan() {
		idx++
		assert.True(t, src_s.Scan())
		assert.Equal(t, expected_s.Bytes(), src_s.Bytes(), "line "+strconv.Itoa(idx))
	}
}
