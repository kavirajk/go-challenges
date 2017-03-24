package drum

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	p := &Pattern{}
	header := [6]byte{}
	if err := binary.Read(f, binary.BigEndian, &header); err != nil {
		return nil, err
	}
	p.Header = fmt.Sprintf("%s", header)

	if err := binary.Read(f, binary.BigEndian, &p.Size); err != nil {
		return nil, err
	}

	lf := &io.LimitedReader{R: f, N: p.Size}

	var version [32]byte

	if err := binary.Read(lf, binary.BigEndian, &version); err != nil {
		return nil, err
	}

	p.Version = strings.TrimRight(fmt.Sprintf("%s", version[:]), "\x00")
	if err := binary.Read(lf, binary.LittleEndian, &p.Tempo); err != nil {
		return nil, err
	}

	// Reading tracks
	tracks := make([]Track, 0)
	for {
		t := Track{}
		err := binary.Read(lf, binary.BigEndian, &t.ID)
		if err == io.EOF {
			break
		}
		var length int32
		if err := binary.Read(lf, binary.BigEndian, &length); err != nil {
			return nil, err
		}
		title := make([]byte, length)

		if err := binary.Read(lf, binary.BigEndian, &title); err != nil {
			return nil, err
		}
		t.Name = fmt.Sprintf("%s", title)
		if err := binary.Read(lf, binary.BigEndian, &t.Steps); err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	p.Tracks = tracks

	return p, nil
}
