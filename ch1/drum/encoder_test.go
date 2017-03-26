package drum

import (
	"path"
	"reflect"
	"testing"
)

func TestEncodePattern(t *testing.T) {
	tests := []struct {
		ipath, opath string
	}{
		{"pattern_1.splice", "pattern_1_o.splice"},
		{"pattern_2.splice", "pattern_2_o.splice"},
		{"pattern_3.splice", "pattern_3_o.splice"},
		{"pattern_4.splice", "pattern_4_o.splice"},
		{"pattern_5.splice", "pattern_5_o.splice"},
	}

	for _, test := range tests {
		expPattern, err := DecodeFile(path.Join("fixtures", test.ipath))
		if err != nil {
			t.Fatalf("error in decoding the pattern %s - %v", test.ipath, err)
		}
		if err := EncodePattern(expPattern, path.Join("fixtures", test.opath)); err != nil {
			t.Fatalf("error in encoding pattern to file %s - %v", test.opath, err)
		}

		outPattern, err := DecodeFile(path.Join("fixtures", test.opath))
		if err != nil {
			t.Fatalf("error in decoding the actual pattern %s - %v", test.opath, err)
		}

		if !reflect.DeepEqual(outPattern, expPattern) {
			t.Fatalf("pattern mismatch. Seriously wrong with encoding pattern %s - %s", outPattern, expPattern)
		}
	}
}
