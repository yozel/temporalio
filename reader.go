package temporalio

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type Reader struct {
	storage io.Reader
	start   time.Time
}

func NewReader(storage io.Reader) *Reader {
	return &Reader{storage: storage, start: time.Time{}}
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.start.IsZero() {
		r.start = time.Now()
	}
	duration, _n, err := r.read()
	if err != nil {
		return 0, err
	}
	n := int(_n)
	if len(p) < n {
		return 0, fmt.Errorf("buffer is too small: %d < %d", len(p), _n)
	}

	sleepDuration := time.Until(r.start.Add(duration))
	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
	}

	buffer := make([]byte, n)
	var _n2 int
	_n2, err = r.storage.Read(buffer)
	if err != nil {
		return 0, err
	}
	if _n2 != n {
		return 0, fmt.Errorf("read %d bytes, but expected %d bytes", _n, n)
	}

	return copy(p, buffer), nil
}

func (r *Reader) read() (duration time.Duration, n int64, err error) {
	err = binary.Read(r.storage, binary.LittleEndian, &duration)
	if err != nil {
		if err == io.EOF {
			return 0, 0, err
		}
		return 0, 0, fmt.Errorf("read duration: %s", err)
	}
	err = binary.Read(r.storage, binary.LittleEndian, &n)
	if err != nil {
		return 0, 0, fmt.Errorf("read n: %s", err)
	}
	return duration, n, nil
}
