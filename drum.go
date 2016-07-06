package main
// package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information

import (
	"fmt"
	"strconv"
)

// Steps is a representation of the 16 bytes defining a drum machine pattern.
type Steps [16]byte

// String representation of the Steps type; quite complivated format...
func (st Steps) String() string {
	s := ""
	for cnt, step := range st {

		if cnt%4 == 0 {
			s = fmt.Sprintf("%s|", s)
		}
		if step == 0x01 {
			s = fmt.Sprintf("%sx", s)
		} else {
			s = fmt.Sprintf("%s-", s)
		}
	}
	s = fmt.Sprintf("%s|", s)
	return s
}

// Track is a definition of a single drum machine track
type Track struct {

	// ID of the pattern
	ID uint32

	// Name of the pattern, using byte instead of string to simplify the decoding
	Name string

	// exactly 16 steps: only values 0 or 1 are allowed
	Steps Steps
}

// String representation of the drum machine track.
func (t *Track) String() string { return fmt.Sprintf("(%d) %s\t%s", t.ID, t.Name, t.Steps.String()) }

// NewTrack creates a new empty instance of Track type.
func NewTrack(id uint32, name string) *Track {
	var s Steps
	return &Track{id, name, s}
}

// SetSteps sets the needed steps (with given varags indexes) in Track.
func (t *Track) SetSteps(steps ...int) error {

	for s := range steps {
		if s < 0 || s > 16 {
			return fmt.Errorf("Step index (%d) out of range (valid: 1 - 16).\n", s)
		}
		steps[s-1] = 0x01 // use hex to emphasize this is byte
	}
	return nil
}

/*
//
func (t *Track) Read(b []byte) error {

     length := len(b)
     if length < 22 {
         return t, fmt.Errorf("Buffer is too short to decode the track.\n")
     }
     // TODO: implement an io.Reader for track decoding

    return t, nil
}
*/

// Pattern is the high level representation of the drum pattern contained in a .splice file.
type Pattern struct {

	// the name of the source file
	Filename string

	// decoded HW version string
	Version string

	// decode tempo
	Tempo float32

	// A list of decode drum machine tracks
	tracks []Track
}

// NewPattern creates a new empty instance of Pattern type.
func NewPattern() *Pattern {
	return &Pattern{"", "", 0.0, make([]Track, 0)}
}

// AddTrack appneds an additional track to the Drum machine pattern.
func (p *Pattern) AddTrack(tr *Track) { p.tracks = append(p.tracks, *tr) }

// String representation of the Pattern.
func (p *Pattern) String() string {

	s := fmt.Sprintf("%s\nSaved with HW Version: %s\nTempo: %s\n",
		p.Filename, p.Version, strconv.FormatFloat(float64(p.Tempo), 'f', -1, 32))
	for _, track := range p.tracks {
		s = fmt.Sprintf("%s%s\n", s, track.String())
	}
	return s
}

// GetTrack returns a track according to given index.
func (p *Pattern) GetTrack(index int) (*Track, error) {

	if index < 0 || index >= len(p.tracks) {
		return nil, fmt.Errorf("Track index (%d) out of range.\n", index)
	}
	return &p.tracks[index], nil
}

// GetTrackByName returns a track according to given name. If not found, nil is returned.
func (p *Pattern) GetTrackByName(name string) *Track {

	for _, t := range p.tracks {
		if name == string(t.Name) {
			return &t
		}
	}
	return nil

}
