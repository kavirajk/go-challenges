package drum

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
// TODO: implement
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

	var version [11]byte

	binary.Read(f, binary.BigEndian, &version)

	var i int
	for i = 0; i < 11; i++ {
		if version[i] == 0x00 {

			break
		}
	}
	p.Version = fmt.Sprintf("%s", version[:i])
	dummy = make([]byte, 21)

	binary.Read(f, binary.BigEndian, &dummy)
	binary.Read(f, binary.LittleEndian, &p.Tempo)

	// Reading tracks
	tracks := make([]Track, 0)
	for {
		t := Track{}
		err := binary.Read(f, binary.BigEndian, &t.ID)
		if err == io.EOF {
			break
		}
		dummy = make([]byte, 3)
		binary.Read(f, binary.BigEndian, &dummy)
		var length byte
		binary.Read(f, binary.BigEndian, &length)
		title := make([]byte, length)
		binary.Read(f, binary.BigEndian, &title)
		t.Title = fmt.Sprintf("%s", title)
		binary.Read(f, binary.BigEndian, &t.States)
		tracks = append(tracks, t)
	}
	p.Tracks = tracks

	return p, nil
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
// TODO: implement
type Pattern struct {
	Header  [6]byte
	Size    byte
	Version string
	Tempo   float32
	Tracks  []Track
}

func (p Pattern) String() string {
	var ret string
	ret = fmt.Sprintf("Saved with HW Version: %s\n", p.Version)
	if (p.Tempo - float32(int(p.Tempo))) == 0 {
		ret += fmt.Sprintf("Tempo: %.0f\n", p.Tempo)
	} else {
		ret += fmt.Sprintf("Tempo: %.1f\n", p.Tempo)
	}
	for i := range p.Tracks {
		ret += fmt.Sprint(p.Tracks[i])
		ret += "\n"
	}
	return ret
}

type Track struct {
	ID     uint8
	Title  string
	States [16]byte
}

func (t Track) String() string {
	var ret string
	ret = fmt.Sprintf("(%d) %s\t", t.ID, t.Title)
	ret += "|"

	for i := range t.States {
		if t.States[i] == 1 {
			ret += fmt.Sprintf("%c", 'x')
		} else {
			ret += fmt.Sprintf("%c", '-')
		}
		if (i+1)%4 == 0 {
			ret += "|"
		}
	}
	return ret
}
