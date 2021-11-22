package ptest

import (
	"fmt"
	"testing"
)

var pairs = []struct {
	k string
	v string
}{
	{"polaris", " Andy "},
	{"studygolang", "Go Language "},
	{"stdlib", "Go Lib "},
	{"polaris1", " Andy 1"},
	{"studygolang1", "Go Language 1"},
	{"stdlib1", "Go Lib 1"},
	{"polaris2", " Andy 2"},
	{"studygolang2", "Go Language 2"},
	{"stdlib2", "Go Lib 2"},
	{"polaris3", " Andy 3"},
	{"studygolang3", "Go Language 3"},
	{"stdlib3", "Go Lib 3"},
	{"polaris4", " Andy 4"},
	{"studygolang4", "Go Language 4"},
	{"stdlib4", "Go Lib 4"},
}

// TestWriteToMap need to run before TestReadFromMap
func TestWriteToMap(t *testing.T) {
	t.Parallel()
	for _, tt := range pairs {
		WriteToMap(tt.k, tt.v)
	}
}

func TestReadFromMap(t *testing.T) {
	t.Parallel()
	for _, tt := range pairs {
		actual := ReadFromMap(tt.k)
		if actual != tt.v {
			t.Errorf("the value of key(%s) is [%s], expected: %s", tt.k, actual, tt.v)
		}
	}
}

func Benchmark_fmt_Sprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("hello ")
	}
}

func Sum(numbers []int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

func TestSum(t *testing.T) {
	tc := []struct {
		value    []int
		expected int
	}{
		{[]int{2, 2, 2, 4}, 10},
		{[]int{-1, -2, -3, -4, 5}, -5},
		{[]int{-1, -2, 0, 1, 2, 3, 4}, 7},
	}
	for _, n := range tc {
		got := Sum(n.value)
		if got != n.expected {
			t.Errorf("expect %d, got %d\n", n.expected, got)
		}
	}
}

/*
func TestTime(t *testing.T) {
	testCases := []struct {
		gmt  string
		loc  string
		want string
	}{
		{"12:31", "Europe/Zurich", "13:05"}, // incorrect location name
		{"12:31", "America/New_York", "07:34"},
		{"08:08", "Australia/Sydney", "18:12"},
	}

	for _, tc := range testCases {
		loc, err := time.LoadLocation(tc.loc)
		if err != nil {
			t.Fatalf("could not load location %q", tc.loc)
		}
		gmt, _ := time.Parse("15:04", tc.gmt)
		if got := gmt.In(loc).Format("15:04"); got != tc.want {
			t.Errorf("In(%s, %s) = %s; want %s", tc.gmt, tc.loc, got, tc.want)
		}
	}
}
*/
