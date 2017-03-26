// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"bytes"
	"fmt"
)

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
	w := &bytes.Buffer{}
	fmt.Fprintf(w, "Saved with HW Version: %s\n", p.Version)
	fmt.Fprintf(w, "Tempo: %g\n", p.Tempo)
	for i := range p.Tracks {
		w.WriteString(fmt.Sprint(p.Tracks[i]))
		w.WriteByte('\n')
	}
	return w.String()
}

func (p *Pattern) AddCowbell(cowbell [16]byte) {
	var maxID uint8
	for i := range p.Tracks {
		if p.Tracks[i].ID > maxID {
			maxID = p.Tracks[i].ID
		}
		if p.Tracks[i].Name == "cowbell" {
			p.Tracks[i].Steps = cowbell
			return
		}
	}
	// cowbell not found. So add it
	p.Tracks = append(p.Tracks, Track{
		ID:    maxID + 1,
		Name:  "cowbell",
		Steps: cowbell,
	})

	// increase the size
	p.Size += 8 + 7 + 16 // ID + name + steps
}

type Track struct {
	ID    uint8
	Name  string
	Steps [16]byte
}

func (t Track) String() string {
	w := &bytes.Buffer{}
	fmt.Fprintf(w, "(%d) %s\t", t.ID, t.Name)
	for i := range t.Steps {
		if i%4 == 0 {
			w.WriteByte('|')
		}

		if t.Steps[i] == 1 {
			w.WriteByte('x')
		} else {
			w.WriteByte('-')
		}
	}
	w.WriteByte('|')
	return w.String()
}
