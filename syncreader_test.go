package syncreader

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/bmizerany/assert"
)

type ProgressReader struct {
	Current, Total, Expected int64
	Progress                 func(current, total, expected int64)
	isFinished               bool
}

func (p *ProgressReader) Read(b []byte) (n int, err error) {
	if p.isFinished {
		return 0, io.EOF
	}
	n = len(b)
	p.isFinished = p.calculate(int64(n))
	p.Progress(p.Current, p.Total, p.Expected)
	return
}

func (p *ProgressReader) calculate(i int64) bool {
	p.Current += i
	if p.Current > p.Total {
		p.Current = p.Total
	}
	p.Expected = p.Total - p.Current
	return p.Current == p.Total
}

func TestNew(t *testing.T) {
	filename := "syncreader_test.go"
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}

	fs, err := os.Stat(filename)
	if err != nil {
		log.Fatalln(err)
	}

	p := new(ProgressReader)
	p.Total = fs.Size()
	p.Progress = func(c, t, e int64) {
		log.Println(c, t, e)
	}
	b := new(bytes.Buffer)
	r := New(f, p)
	_, err = b.ReadFrom(r)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, fs.Size(), int64(b.Len()))
}
