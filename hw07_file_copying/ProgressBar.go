package main

import (
	"fmt"
	"strings"
)

type ProgressBar interface {
	Total() int64
	Update(proc int64)
	Render()
}

type TxtProgressBar struct {
	total, processed int64
	barWidth         int

	bar string
}

func NewTxtProgressBar(total int64, barWidth int) *TxtProgressBar {
	return &TxtProgressBar{total: total, barWidth: barWidth}
}

func (t TxtProgressBar) Total() int64 {
	return t.total
}

func (t *TxtProgressBar) Update(proc int64) {
	t.processed = proc
}

func (t *TxtProgressBar) Render() {
	percent := float64(t.processed) / float64(t.total) * 100
	filled := int(float64(t.barWidth) * float64(t.processed) / float64(t.total))

	bar := strings.Repeat("#", filled) + strings.Repeat(" ", t.barWidth-filled)

	t.bar = fmt.Sprintf("\r[%s] %.1f%% (%d/%d)", bar, percent, t.processed, t.total)
	fmt.Print(t.bar)
}
