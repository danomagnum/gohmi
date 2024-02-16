package main

import "slices"

type ReadLength struct {
	Start int
	Last  int
}

func (r ReadLength) Length() int {
	return r.Last - r.Start + 1
}

// The purpose of this function is to optimize / minimize the number of reads for (offset, count) type reads
// that modbus, s7, and other protocols use.
func compute_reads(addresses []int, factor int) []ReadLength {

	slices.Sort(addresses)

	results := make([]ReadLength, 0, len(addresses))
	var result ReadLength

	result = ReadLength{Start: addresses[0], Last: addresses[0]}
	for i := range addresses {
		addr := addresses[i]
		delta := addr - result.Last
		if delta > factor {
			results = append(results, result)
			result = ReadLength{Start: addr, Last: addr}
			continue
		}
		result.Last = addr
	}
	results = append(results, result)

	return results
}
