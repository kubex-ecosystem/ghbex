// Package state provides types and functions for managing application state.
package state

import "sync/atomic"

type Stage uint64

const (
	StageRunsCleanup Stage = 1 << iota
	StageArtifactsCleanup
	StageReleaseCleanup
	StageNotify
	StageReportPersist
)

type FlagSet struct{ v atomic.Uint64 }

func (f *FlagSet) Enable(s Stage) { f.v.Add(uint64(s)) }
func (f *FlagSet) Disable(s Stage) {
	for {
		old := f.v.Load()
		if f.v.CompareAndSwap(old, old&^uint64(s)) {
			return
		}
	}
}
func (f *FlagSet) Has(s Stage) bool { return f.v.Load()&uint64(s) != 0 }
