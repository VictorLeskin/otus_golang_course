package hw07_file_copying

type ProgressBar struct {
	processed, expected int64
}

func (t ProgressBar) len() int64 {
	return t.expected
}

func (t *ProgressBar) update(proc int64) {
	t.processed = proc
}

func NewProgressBar(expected int64) *ProgressBar {
	return &ProgressBar{expected: expected}
}
