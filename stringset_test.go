// Copyright 2016 John Jeffery <john@jeffery.id.au>. All rights reserved.

package stringset

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		set1  Set
		set2  Set
		equal bool
	}{
		{
			set1:  New("a", "b"),
			set2:  New("b", "a"),
			equal: true,
		},
		{
			set1:  New("a", "b", "a"),
			set2:  New("b", "a"),
			equal: true,
		},
		{
			set1:  nil,
			set2:  New(),
			equal: true,
		},
		{
			set1:  nil,
			set2:  New("x"),
			equal: false,
		},
		{
			set1:  New("X"),
			set2:  New("x"),
			equal: false,
		},
	}
	for i, tt := range tests {
		got := tt.set1.Equal(tt.set2)
		if want := tt.equal; got != want {
			t.Errorf("%d: expected=%v, actual=%v", i, want, got)
		}
		got = tt.set2.Equal(tt.set1)
		if want := tt.equal; got != want {
			t.Errorf("%d: expected=%v, actual=%v", i, want, got)
		}
	}
}

func TestAdd(t *testing.T) {
	var set Set
	set = Add(set, "A")
	if got, want := set.Len(), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	set = Add(set, "B")
	if got, want := set.Len(), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	set.Add("C", "D", "E")
	if got, want := set.Len(), 5; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	// Add method on nil set.
	set = nil
	set.Add("A")
	if got, want := set.Len(), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		set    Set
		remove []string
		want   Set
	}{
		{
			set:    New("A", "B", "C"),
			remove: []string{"Z"},
			want:   New("A", "B", "C"),
		},
		{
			set:    New("A", "B", "C"),
			remove: []string{"A"},
			want:   New("B", "C"),
		},
		{
			set:    New("A", "B", "C"),
			remove: []string{"A", "C"},
			want:   New("B"),
		},
		{
			set:    nil,
			remove: []string{"A"},
			want:   nil,
		},
	}
	for i, tt := range tests {
		tt.set.Remove(tt.remove...)
		if !tt.set.Equal(tt.want) {
			t.Errorf("%d: want %v, got %v", i, tt.want, tt.set)
		}
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		set  Set
		want int
	}{
		{
			set:  nil,
			want: 0,
		},
		{
			set:  New(),
			want: 0,
		},
		{
			set:  New("1"),
			want: 1,
		},
		{
			set:  New("1", "2", "3"),
			want: 3,
		},
	}
	for i, tt := range tests {
		if got := tt.set.Len(); got != tt.want {
			t.Errorf("%d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		set  Set
		s    string
		want bool
	}{
		{
			set:  nil,
			s:    "A",
			want: false,
		},
		{
			set:  New(),
			s:    "B",
			want: false,
		},
		{
			set:  New("1"),
			s:    "1",
			want: true,
		},
		{
			set:  New("1a", "2b", "3c"),
			s:    "2A",
			want: false,
		},
		{
			set:  New("1a", "2b", "3c"),
			s:    "2b",
			want: true,
		},
	}
	for i, tt := range tests {
		if got := tt.set.Contains(tt.s); got != tt.want {
			t.Errorf("%d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		set  Set
		want []string
	}{
		{
			set:  nil,
			want: nil,
		},
		{
			set:  New(),
			want: nil,
		},

		{
			set:  New("1"),
			want: []string{"1"},
		},
		{
			set:  New("1a", "2b", "3c"),
			want: []string{"1a", "2b", "3c"},
		},
		{
			set:  New("2b", "1a", "3c"),
			want: []string{"1a", "2b", "3c"},
		},
	}
	for i, tt := range tests {
		if got := tt.set.Values(); !compareStringSlices(got, tt.want) {
			t.Errorf("%d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		set   Set
		want  string
		want2 string
	}{
		{
			set:   nil,
			want:  "[]",
			want2: "stringset.Set(nil)",
		},
		{
			set:   New(),
			want:  "[]",
			want2: "stringset.Set(nil)",
		},

		{
			set:   New("1"),
			want:  `[1]`,
			want2: `stringset.Set{"1"}`,
		},
		{
			set:   New("1a", "2b", "3c"),
			want:  `[1a 2b 3c]`,
			want2: `stringset.Set{"1a", "2b", "3c"}`,
		},
		{
			set:   New("2b", "1a", "3c"),
			want:  `[1a 2b 3c]`,
			want2: `stringset.Set{"1a", "2b", "3c"}`,
		},
	}
	for i, tt := range tests {
		if got := fmt.Sprintf("%v", tt.set); got != tt.want {
			t.Errorf("%d: got %v, want %v", i, got, tt.want)
		}
		if got := tt.set.String(); got != tt.want {
			t.Errorf("%d: got %v, want %v", i, got, tt.want)
		}
		if got := fmt.Sprintf("%#v", tt.set); got != tt.want2 {
			t.Errorf("%d: got %v, want %v", i, got, tt.want2)
		}
		if got := tt.set.GoString(); got != tt.want2 {
			t.Errorf("%d: got %v, want %v", i, got, tt.want2)
		}
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		set  Set
		sep  string
		want string
	}{
		{
			set:  nil,
			sep:  ",",
			want: "",
		},
		{
			set:  New(),
			sep:  ".",
			want: "",
		},
		{
			set:  New("1"),
			sep:  "-",
			want: "1",
		},
		{
			set:  New("1a", "2b", "3c"),
			sep:  ", ",
			want: "1a, 2b, 3c",
		},
		{
			set:  New("2b", "1a", "3c"),
			sep:  ".",
			want: "1a.2b.3c",
		},
	}
	for i, tt := range tests {
		if got := tt.set.Join(tt.sep); got != tt.want {
			t.Errorf("%d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		set  Set
		want string
	}{
		{
			set:  nil,
			want: "null",
		},
		{
			set:  New(),
			want: "null",
		},
		{
			set:  New("1"),
			want: `["1"]`,
		},
		{
			set:  New("1a", "2b", "3c"),
			want: `["1a","2b","3c"]`,
		},
		{
			set:  New("2b", "1a", "3c"),
			want: `["1a","2b","3c"]`,
		},
	}
	for i, tt := range tests {
		data, err := json.Marshal(tt.set)
		if err != nil {
			t.Error(err)
			continue
		}
		if got := string(data); got != tt.want {
			t.Errorf("%d: got %v, want %v", i, got, tt.want)
			continue
		}
		var set Set
		if err := json.Unmarshal(data, &set); err != nil {
			t.Error(err)
			continue
		}
		if !set.Equal(tt.set) {
			t.Errorf("%d: got %v, want %v", i, set, tt.set)
		}
	}

	// Test invalid JSON input.
	{
		var set Set
		data := []byte(`["1","2",3]`)
		if err := json.Unmarshal(data, &set); err == nil {
			t.Errorf("got nil, expected error")
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		set  Set
		want string
	}{
		{
			set:  nil,
			want: "[]",
		},
		{
			set:  New(),
			want: "[]",
		},

		{
			set:  New("1"),
			want: "[1]",
		},
		{
			set:  New("1a", "2b", "3c"),
			want: "[1a 2b 3c]",
		},
		{
			set:  New("2b", "1a", "3c"),
			want: "[1a 2b 3c]",
		},
	}
	for i, tt := range tests {
		if got, want := tt.set.String(), tt.want; got != want {
			t.Errorf("%d: got %v, want %v", i, got, want)
		}
	}
}

func compareStringSlices(v1, v2 []string) bool {
	if v1 == nil && v2 == nil {
		return true
	} else if v1 == nil {
		return false
	} else if v2 == nil {
		return false
	}
	if len(v1) != len(v2) {
		return false
	}
	for i, s1 := range v1 {
		if s1 != v2[i] {
			return false
		}
	}
	return true
}
