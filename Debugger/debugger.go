package Debugger

import (
	"fmt"
	"strings"

	wt "github.com/PlayerR9/LyneDbg/Wait"
)

type DebugStringer interface {
	DebugString() string
}

type DebugObserver[T DebugStringer] struct {
	data T

	waitee *wt.Waitee
}

func (d *DebugObserver[T]) Notify(change ChangeMsg) {
	var lines []string

	lines = append(lines, "Change detected:")
	lines = append(lines, change.String())
	lines = append(lines, "")
	lines = append(lines, "Debug Information:")
	lines = append(lines, d.data.DebugString())

	fmt.Println(strings.Join(lines, "\n"))

	d.waitee.Wait()
}

func NewDebugObserver[T DebugStringer](data T) *DebugObserver[T] {
	return &DebugObserver[T]{
		data: data,
	}
}

type Process struct {
	wait *wt.Waitee

	// other fields
}

func (p *Process) MyFunction() {
	// do something

	p.wait.Wait()

	// continue
}
