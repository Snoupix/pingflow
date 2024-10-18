package utils

import "sync/atomic"

const MAX_IDX = 0xFF // Between 1 and 255 seems alright since it shouldn't have that much jobs overlapping
const MIN_IDX = 0x1

var work_id_idx atomic.Uint32

func Get() uint32 {
    return work_id_idx.Load()
}

// Increments and returns the current job index
func NewWorkIdx() uint32 {
    var new_val uint32

    for {
        val := work_id_idx.Load()
        new_val = (val + 1) % MAX_IDX
        if new_val < MIN_IDX {
            new_val = MIN_IDX
        }

        if work_id_idx.CompareAndSwap(val, new_val) {
            break
        }
    }

    return new_val
}
