package main

import "testing"

func TestBulkCombine(t *testing.T) {
	tests := []struct {
		name   string
		addr   []int
		factor int
		want   []ReadLength
	}{
		{
			"Test0",
			[]int{1, 2, 6, 7, 8},
			3,
			[]ReadLength{{Start: 1, Last: 2}, {Start: 6, Last: 8}},
		},
		{
			"Test1",
			[]int{7, 6, 2, 8, 1},
			3,
			[]ReadLength{{Start: 1, Last: 2}, {Start: 6, Last: 8}},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			have := compute_reads(test.addr, test.factor)
			if len(have) != len(test.want) {
				t.Errorf("got wrong number of results for %s: got %d want %d", test.name, len(have), len(test.want))
				return
			}

			for i, w := range test.want {
				h := have[i]
				if h.Last != w.Last {
					t.Errorf("got wrong last for %s, %d: got %d want %d", test.name, i, h.Last, w.Last)
					return
				}
				if h.Start != w.Start {
					t.Errorf("got wrong start for %s, %d: got %d want %d", test.name, i, h.Start, w.Start)
					return
				}

			}
		})

	}
}
