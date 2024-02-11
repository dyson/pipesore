package pipeline

import (
	"io"
	"sync"
)

// NewPipeline returns new pipeline given	a io.Reader.
func NewPipeline(r io.Reader) *pipeline {
	return &pipeline{r: r}
}

// A pipeline contains the io.Reader the next filter is to read from as well as
// a mutex protected error for all filter errors to be written to.
type pipeline struct {
	r io.Reader

	err error
	mu  sync.Mutex
}

// SetError writes an error to the pipline.
func (p *pipeline) SetError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.err = err
}

// Filter takes a filter function and filters the current pipelines io.Reader
// with it. It writes the output to an io.Pipe() and sets the pipelines
// io.Reader to the io.Pipe io.Reader ready for the next filter to consume. If
// the filter errors the error is set on the pipeline.
func (p *pipeline) Filter(filter func(io.Reader, io.Writer) error) {
	r := p.r

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		err := filter(r, pw)
		p.SetError(err)
	}()

	p.r = pr
}

// Output copies from the pipelines io.Reader and writes it to the provided
// io.Writer.
func (p *pipeline) Output(out io.Writer) (int64, error) {
	i, err := io.Copy(out, p.r)
	if p.err != nil {
		return i, p.err
	}

	return i, err
}
