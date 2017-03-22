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
	binary.Read(f, binary.BigEndian, &p.Header)

	dummy := make([]byte, 7)
	binary.Read(f, binary.BigEndian, &dummy)

	binary.Read(f, binary.BigEndian, &p.Size)

	lf := &io.LimitedReader{R: f, N: int64(p.Size)}

	var version [32]byte

	binary.Read(lf, binary.BigEndian, &version)

	p.Version = strings.TrimRight(fmt.Sprintf("%s", version[:]), "\x00")
	binary.Read(lf, binary.LittleEndian, &p.Tempo)

	// Reading tracks
	tracks := make([]Track, 0)
	for {
		t := Track{}
		err := binary.Read(lf, binary.BigEndian, &t.ID)
		if err == io.EOF {
			break
		}
		var length int32
		binary.Read(lf, binary.BigEndian, &length)
		title := make([]byte, length)
		binary.Read(lf, binary.BigEndian, &title)
		t.Name = fmt.Sprintf("%s", title)
		binary.Read(lf, binary.BigEndian, &t.Steps)
		tracks = append(tracks, t)
	}
	p.Tracks = tracks

	return p, nil
}
