package main

import (
	"fmt"
	"github.com/cespare/xxhash/v2"
	"iter"
	"os"
	"runtime"
	"runtime/pprof"
	"syscall"
	"unsafe"
)

type StationData struct {
	id      uint64
	name    string
	minimum int16
	maximum int16
	sum     int64
	count   int
}

type ProcessedResults struct {
	items [32768]StationData
}

func (p *ProcessedResults) get(id uint64) (*StationData, *StationData) {
	index := uint16(id) >> 1
	if p.items[index].count == 0 {
		return nil, &p.items[index]
	}
	for {
		if p.items[index].id == id {
			return &p.items[index], nil
		}
		index = (index + 1) % 32768
		if p.items[index].count == 0 {
			return nil, &p.items[index]
		}
	}
}

func (p *ProcessedResults) entries() iter.Seq[*StationData] {
	return func(yield func(*StationData) bool) {
		for i := range p.items {
			if p.items[i].count != 0 {
				if !yield(&p.items[i]) {
					return
				}
			}
		}
	}
}

func detectDelimiter(data *byte) uint64 {
	n := *(*uint64)(unsafe.Pointer(data)) ^ (';' * 0x0101010101010101)
	return (n - 0x0101010101010101) &^ n & 0x8080808080808080
}

func main() {
	inputFile := "measurements_1b.txt"
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}

	if os.Getenv("PROFILE") != "" {
		profFileName := os.Args[0] + ".prof"
		fmt.Fprintln(os.Stderr, "### profiling!")
		pfile, _ := os.Create(profFileName)
		defer pfile.Close()
		pprof.StartCPUProfile(pfile)
		defer func() {
			pprof.StopCPUProfile()
		}()
	}
	decimalLookup := [65536]int16{}
	for i := range 1000 {
		s := fmt.Sprintf("%d.%d\n", i/10, i%10)
		b := []byte(s)
		asInt := (*uint32)(unsafe.Pointer(&b[0]))
		coolBits := *asInt & 0x0f0f0f0f
		folded := uint16(coolBits) | uint16(coolBits>>12)
		skip := int16(len(s) << 10)
		decimalLookup[folded] = int16(i) | skip
	}


datas := partitionData(mmapFile(inputFile))
results := ProcessedResults{}
resultsCh := make(chan *ProcessedResults)
for _, data := range datas {
        go func() {
                results := ProcessedResults{}
                parseLines(&results, data, &decimalLookup)
                resultsCh <- &results
        }()
}
for range datas {
        q := <-resultsCh
        for i := range q.items {
                if q.items[i].count == 0 {
                        continue
                }
                if rItem, newItem := results.get(q.items[i].id); newItem != nil {
                        *newItem = q.items[i]
                } else {
                        rItem.count += q.items[i].count
                        rItem.sum += q.items[i].sum
                        rItem.minimum = min(rItem.minimum, q.items[i].minimum)
                        rItem.maximum = max(rItem.maximum, q.items[i].maximum)
                }
	}
}
	for v := range results.entries() {
		fmt.Printf(
			"%s;%.1f;%.1f;%.1f\n",
			v.name,
			0.1*float32(v.minimum),
			0.1*float64(v.sum)/float64(v.count),
			0.1*float32(v.maximum),
		)
	}
}

func parseLines(results *ProcessedResults, data []byte, lookup *[65536]int16) {
	pos := 0
	for pos < len(data) {
		lineStart := pos
	skipToDelim:
		for {
			skip := detectDelimiter(&data[pos])
			switch skip & 0xffffffff {
			case 0x80000000:
				pos += 3
				break skipToDelim
			case 0x00800000:
				pos += 2
				break skipToDelim
			case 0x00008000:
				pos += 1
				break skipToDelim
			case 0x00000080:
				break skipToDelim
			}
			switch skip >> 32 {
			case 0x80000000:
				pos += 7
				break skipToDelim
			case 0x00800000:
				pos += 6
				break skipToDelim
			case 0x00008000:
				pos += 5
				break skipToDelim
			case 0x00000080:
				pos += 4
				break skipToDelim
			}
			pos += 8
		}
		split := pos
		stationName := data[lineStart:split]
		id := xxhash.Sum64(stationName)
		pos += 1

		var measurement int16
		if data[pos] == '-' {
			pos += 1
			asInt := (*uint32)(unsafe.Pointer(&data[pos]))
			coolBits := *asInt & 0x0f0f0f0f
			folded := uint16(coolBits) | uint16(coolBits>>12)
			num := lookup[folded]
			measurement = -(num & 0x3ff)
			pos += int(num >> 10)
		} else {
			asInt := (*uint32)(unsafe.Pointer(&data[pos]))
			coolBits := *asInt & 0x0f0f0f0f
			folded := uint16(coolBits) | uint16(coolBits>>12)
			num := lookup[folded]
			measurement = num & 0x3ff
			pos += int(num >> 10)
		}

		item, newItem := results.get(id)
		if newItem != nil {
			newItem.name = string(stationName)
			newItem.id = id
			newItem.minimum = measurement
			newItem.maximum = measurement
			newItem.sum = int64(measurement)
			newItem.count = 1
		} else {
			item.maximum = max(item.maximum, measurement)
			item.minimum = min(item.minimum, measurement)
			item.sum += int64(measurement)
			item.count += 1
		}
	}
}

func mmapFile(filename string) []byte {
	f, _ := os.Open(filename)
	defer f.Close()
	fi, _ := f.Stat()
	size := fi.Size()
	data, err := syscall.Mmap(
		int(f.Fd()),
		0,
		int(size)+8,
		syscall.PROT_READ,
		syscall.MAP_SHARED|syscall.MAP_POPULATE,
	)
	if err != nil {
		panic(err)
	}
	syscall.Madvise(data, syscall.MADV_SEQUENTIAL)
	return data[0 : len(data)-8]
}

func partitionData(data []byte) [][]byte {
        n := runtime.NumCPU()
        partitions := make([][]byte, n)
        partitionSize := len(data) / n
        prevEnd := 0
        for i := range n {
                start := prevEnd
                end := max(start, partitionSize*(i+1))
                if i == n - 1 {
                        end = len(data)
                } else {
                        for data[end-1] != '\n' {
                                end += 1
                        }
                }
                prevEnd = end
                partitions[i] = data[start:end]
        }
        return partitions
}

