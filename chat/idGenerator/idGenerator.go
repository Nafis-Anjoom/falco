package idGenerator

import (
	"errors"
	"sync"
	"time"
)

const (
	timestampBits      = 41
	machineIdBits      = 10
	SequenceNumberBits = 12

	maxMachineId = -1 ^ (-1 << machineIdBits)
	maxSequence  = -1 ^ (-1 << SequenceNumberBits)
	maxTimestamp = -1 ^ (-1 << timestampBits)

	timestampBitShift = machineIdBits + SequenceNumberBits

	// 2024-11-01 00:00:00 +0000 UTC
	epoch = 1730419200000
)

var (
    invalidMachineId = errors.New("MachineId out of range. Max 1023, Min 0")
)

type IdGenerator struct {
	mutex         *sync.Mutex
	lastTimestamp int64
	machineId     int64
	sequence      int64
}

func NewIdGenerator(machineId int64) (*IdGenerator, error) {
    if machineId < 0 || machineId > maxMachineId {
        return nil, invalidMachineId
    }
	return &IdGenerator{machineId: machineId}, nil
}

func (gen *IdGenerator) Generate() int64 {
    gen.mutex.Lock()
    defer gen.mutex.Unlock()

	timestamp := time.Now().UnixMilli() - epoch
	if timestamp == gen.lastTimestamp {
		gen.sequence = (gen.sequence + 1) & maxSequence

		if gen.sequence == 0 {
			for timestamp <= gen.lastTimestamp {
				timestamp = time.Now().UnixMilli() - epoch
			}
		}
	} else {
		gen.sequence = 0
	}

	gen.lastTimestamp = timestamp

	id := (timestamp << timestampBitShift) |
		(gen.machineId << SequenceNumberBits) |
		gen.sequence

	return id
}
