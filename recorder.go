package temporalio

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type Recorder struct {
	storage io.Writer
	start   time.Time
}

func NewRecorder(storage io.Writer) *Recorder {
	return &Recorder{storage: storage, start: time.Time{}}
}

func (r Recorder) Write(p []byte) (n int, err error) {
	if r.start.IsZero() {
		r.start = time.Now()
	}
	d := time.Since(r.start)
	err = binary.Write(r.storage, binary.LittleEndian, d)
	if err != nil {
		return 0, fmt.Errorf("record error: %v", err)
	}
	err = binary.Write(r.storage, binary.LittleEndian, int64(len(p)))
	if err != nil {
		return 0, fmt.Errorf("record error: %v", err)
	}
	n, err = r.storage.Write(p)
	if err != nil {
		return 0, fmt.Errorf("record error: %v", err)
	}
	if n != len(p) {
		return 0, fmt.Errorf("recorded %d bytes, but expect %d bytes", n, len(p))
	}
	return n, nil
}
