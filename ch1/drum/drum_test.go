package drum

import (
	"fmt"
	"path"
	"testing"
)

func TestAddCowbell(t *testing.T) {
	tests := []struct {
		ipath, opath string
		cowbell      [16]byte
		expPattern   string
	}{
		{
			"pattern_1.splice",
			"pattern_1_o.splice",
			[16]byte{
				0x01, 0x01, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01},
			`Saved with HW Version: 0.808-alpha
Tempo: 120
(0) kick	|x---|x---|x---|x---|
(1) snare	|----|x---|----|x---|
(2) clap	|----|x-x-|----|----|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(4) hh-close	|x---|x---|----|x--x|
(5) cowbell	|xxxx|-x--|---x|xxxx|
`,
		},
		{
			"pattern_2.splice",
			"pattern_2_o.splice",
			[16]byte{
				0x01, 0x01, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01},
			`Saved with HW Version: 0.808-alpha
Tempo: 98.4
(0) kick	|x---|----|x---|----|
(1) snare	|----|x---|----|x---|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(5) cowbell	|xxxx|-x--|---x|xxxx|
`,
		},
		{
			"pattern_3.splice",
			"pattern_3_o.splice",
			[16]byte{
				0x01, 0x01, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01},
			`Saved with HW Version: 0.808-alpha
Tempo: 118
(40) kick	|x---|----|x---|----|
(1) clap	|----|x---|----|x---|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(5) low-tom	|----|---x|----|----|
(12) mid-tom	|----|----|x---|----|
(9) hi-tom	|----|----|-x--|----|
(41) cowbell	|xxxx|-x--|---x|xxxx|
`,
		},
	}

	for _, test := range tests {
		pattern, err := DecodeFile(path.Join("fixtures", test.ipath))
		if err != nil {
			t.Fatalf("error in decoding pattern %s - %v", test.ipath, err)
		}
		pattern.AddCowbell(test.cowbell)
		if err := EncodePattern(pattern, path.Join("fixtures", test.opath)); err != nil {
			t.Fatalf("error in encoding pattern %s - %v", test.opath, err)
		}

		outPattern, err := DecodeFile(path.Join("fixtures", test.opath))
		if err != nil {
			t.Fatalf("error in decoding new pattern %s -%v", test.opath, err)
		}

		if fmt.Sprint(outPattern) != test.expPattern {
			t.Fatalf("something wrong in adding cowbell for got %s - expected %s", fmt.Sprint(outPattern), test.expPattern)
		}
	}
}
