package drum

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path"
)

const (

	// buffer length (in bytes)
	buflen = 1024

	// header length is 14 bytes...
	headerlen = 14

	// absolute minimum length for a single track is 22 bytes: 16 bytes for steps,
	// 4 bytes for id and (1 + 1) bytes for track name
	minTrackLen = 22

	// HW version string length is 32 bytes
	hwVerLen = 32
)

// ReadFile reads a .splice file and return a byte buffer.
func ReadFile(fpath string) (b []byte, err error) {

	b = make([]byte, buflen)

	var fin *os.File
	if fin, err = os.Open(fpath); err != nil {
		return b, fmt.Errorf("Path %q cannot be opened\n", fpath)
	}
	defer fin.Close()

	var nbytes int
	if nbytes, err = fin.Read(b); err != nil {
		return b, fmt.Errorf("Error reading from file %q\n", fpath)
	}
	return b[:nbytes], nil
}

// Decodes the tracks from buffer.
func decodeTracks(buf []byte, p *Pattern) error {

	if buf == nil {
		return fmt.Errorf("Empty buffer.\n")
	}

	buflen := len(buf)
	// if buffer is shorter than MIN, data definitely cannot be decoded properly
	if buflen < minTrackLen {
		return fmt.Errorf("Buffer too short.\n")
	}

	var err error

	csr := 0 // current track cursor
	for csr < buflen {

		track := &Track{}

		// read track ID
		b := bytes.NewBuffer(buf[csr : csr+4])
		if err := binary.Read(b, binary.LittleEndian, &track.ID); err != nil {
			return err
		}

		namelen := int(buf[csr+4])
		stepstart := csr + 5 + namelen
		track.Name = string(buf[csr+5 : stepstart]) // get track name

		// get track steps
		buffer := bytes.NewReader(buf[stepstart : stepstart+16])
		if err = binary.Read(buffer, binary.LittleEndian, &track.Steps); err != nil {
			return fmt.Errorf("Cannot read steps for track: name=%q.\n", track.Name)
		}

		p.AddTrack(track)
		csr += 21 + namelen // this is the exact number of bytes for the current track
	}

	return nil
}

// Decodes the data part of the .splice file.
func decodeData(buf []byte, p *Pattern) error {

	var err error

	// decode the HW version string
	p.Version = string(bytes.TrimRight(buf[:hwVerLen], "\x00"))

	// decode tempo from the file
	b := bytes.NewBuffer(buf[hwVerLen : hwVerLen+4])
	if err := binary.Read(b, binary.LittleEndian, &p.Tempo); err != nil {
		return err
	}

	// finally, decode tracks...
	if err = decodeTracks(buf[hwVerLen+4:], p); err != nil {
		return err
	}

	return nil
}

// Decodes the received data.
func decode(buf []byte, p *Pattern) error {

	// check header
	header := buf[:headerlen]
	if !(string(header[:6]) == "SPLICE") {
		return fmt.Errorf("This is not a .splice file\n")
	}

	// get lengh of the data; this is the one byte before EOH two-byte sequence
	datalen := uint(header[headerlen-1])

	//  decode data part: cut the header bytes from buffer...
	if err := decodeData(buf[headerlen:headerlen+datalen], p); err != nil {
		return err
	}
	return nil
}

// DecodeFile decodes the drum machine file found at the provided path and returns a pointer to a parsed pattern
// which is the entry point to the rest of the data.
func DecodeFile(fpath string) (*Pattern, error) {

	// create a buffer from file
	buf, err := ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	p := &Pattern{}
	p.Filename = path.Base(fpath) // save filename into struct

	err = decode(buf, p)
	return p, err
}
