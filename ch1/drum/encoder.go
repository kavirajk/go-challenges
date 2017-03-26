package drum

import (
	"encoding/binary"
	"os"

	"github.com/pkg/errors"
)

var headerSlice = [6]byte{'S', 'P', 'L', 'I', 'C', 'E'}

func EncodePattern(p *Pattern, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := binary.Write(f, binary.BigEndian, headerSlice); err != nil {
		return errors.Wrap(err, "writing headerslice")
	}

	if err := binary.Write(f, binary.BigEndian, p.Size); err != nil {
		return errors.Wrap(err, "writing pattern size")
	}

	var version [32]byte
	copy(version[:], p.Version)

	if err := binary.Write(f, binary.BigEndian, version); err != nil {
		return errors.Wrap(err, "writing pattern version")
	}

	if err := binary.Write(f, binary.LittleEndian, p.Tempo); err != nil {
		return errors.Wrap(err, "writing pattern tempo")
	}

	for _, t := range p.Tracks {

		if err := binary.Write(f, binary.BigEndian, t.ID); err != nil {
			return errors.Wrapf(err, "writing track ID of %s", t.Name)
		}

		if err := binary.Write(f, binary.BigEndian, int32(len(t.Name))); err != nil {
			return errors.Wrapf(err, "writing track Name's length of %s", t.Name)
		}

		name := []byte(t.Name)

		if err := binary.Write(f, binary.BigEndian, name); err != nil {
			return errors.Wrapf(err, "writing track Name of %s", t.Name)
		}

		if err := binary.Write(f, binary.BigEndian, t.Steps); err != nil {
			return errors.Wrapf(err, "writing track steps of %s", t.Name)
		}

	}

	return nil
}
