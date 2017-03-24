// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import "fmt"

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	Header  string
	Size    int64
	Version string
	Tempo   float32
	Tracks  []Track
}

func (p Pattern) String() string {
	var ret string
	ret = fmt.Sprintf("Saved with HW Version: %s\n", p.Version)
	ret += fmt.Sprintf("Tempo: %g\n", p.Tempo)
	for i := range p.Tracks {
		ret += fmt.Sprint(p.Tracks[i])
		ret += "\n"
	}
	return ret
}

type Track struct {
	ID    uint8
	Name  string
	Steps [16]byte
}

func (t Track) String() string {
	var ret string
	ret = fmt.Sprintf("(%d) %s\t", t.ID, t.Name)
	ret += "|"

	for i := range t.Steps {
		if t.Steps[i] == 1 {
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
