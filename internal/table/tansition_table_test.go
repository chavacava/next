package table

import (
	"fmt"
	"testing"
)

func TestNextListAdd(t *testing.T) {

	tt := []struct {
		s         byte
		wantIndex nextIndex
		wantCount symbolCountType
	}{
		{
			s:         byte(0),
			wantIndex: nextIndex(0),
			wantCount: symbolCountType(1),
		},
		{
			s:         byte(250),
			wantIndex: nextIndex(1),
			wantCount: symbolCountType(1),
		},
		{
			s:         byte(0),
			wantIndex: nextIndex(0),
			wantCount: symbolCountType(2),
		},
		{
			s:         byte(0),
			wantIndex: nextIndex(0),
			wantCount: symbolCountType(3),
		},
		{
			s:         byte(250),
			wantIndex: nextIndex(1),
			wantCount: symbolCountType(2),
		},
		{
			s:         byte(128),
			wantIndex: nextIndex(2),
			wantCount: symbolCountType(1),
		},
	}
	nl := newNextList()

	for _, tc := range tt {
		got := nl.add(tc.s)

		if got != tc.wantIndex {
			fmt.Printf("\n%+v\n", nl)
			t.Fatalf("expected i to be %v, got %d", tc.wantIndex, got)
		}

		next := nl.List[got]
		if next.S != tc.s {
			t.Fatalf("expected byte at %v to be %v, got %v", got, tc.s, next.S)
		}

		if next.Count != tc.wantCount {
			t.Fatalf("expected count to be %v, got %v", tc.wantCount, next.Count)
		}
	}
}
