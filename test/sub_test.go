package ptest

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"testing"
	"time"
)

// https://go.dev/blog/subtests
// https://go.dev/blog/fuzz-beta

func TestMain(m *testing.M) {

	// Setup code
	fmt.Println("Setup")

	exitCode := m.Run()

	// Tear down
	fmt.Println("Tear down")
	os.Exit(exitCode)
}

func BenchmarkTemplateParallel(b *testing.B) {
	templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	b.RunParallel(func(pb *testing.PB) {
		// 每个 goroutine 有属于自己的 bytes.Buffer.
		var buf bytes.Buffer
		for pb.Next() {
			// 循环体在所有 goroutine 中总共执行 b.N 次
			buf.Reset()
			templ.Execute(&buf, "World")
		}
	})
}

func BenchmarkFib3(b *testing.B) { benchmarkFib(3, b) }
func BenchmarkFib5(b *testing.B) { benchmarkFib(5, b) }
func BenchmarkFib7(b *testing.B) { benchmarkFib(7, b) }
func BenchmarkFib9(b *testing.B) { benchmarkFib(9, b) }

func BenchmarkFib(b *testing.B) {
	benchmarks := []struct {
		name  string
		value int
	}{
		{"Fib3", 3},
		{"Fib5", 5},
		{"Fib7", 7},
		{"Fib9", 9},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkFib(bm.value, b)
		})
	}
}

func benchmarkFib(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fib(i)
	}
}

func BenchmarkAppendFloat(b *testing.B) {
	benchmarks := []struct {
		name    string
		float   float64
		fmt     byte
		prec    int
		bitSize int
	}{
		{"Decimal", 33909, 'g', -1, 64},
		{"Float", 339.7784, 'g', -1, 64},
		{"Exp", -5.09e75, 'g', -1, 64},
		{"NegExp", -5.11e-95, 'g', -1, 64},
		{"Big", 123456789123456789123456789, 'g', -1, 64},
	}
	dst := make([]byte, 30)
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				strconv.AppendFloat(dst[:0], bm.float, bm.fmt, bm.prec, bm.bitSize)
			}
		})
	}
}

func TestTime(t *testing.T) {
	testCases := []struct {
		gmt  string
		loc  string
		want string
	}{
		{"12:31", "Europe/Zurich", "13:05"},
		{"12:31", "America/New_York", "07:34"},
		{"08:08", "Australia/Sydney", "18:12"},
	}
	fmt.Println("Setup TT")
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s in %s", tc.gmt, tc.loc), func(t *testing.T) {
			loc, err := time.LoadLocation(tc.loc)
			if err != nil {
				t.Fatal("could not load location")
			}
			gmt, _ := time.Parse("15:04", tc.gmt)
			if got := gmt.In(loc).Format("15:04"); got != tc.want {
				t.Errorf("got %s; want %s", got, tc.want)
			}
		})
	}
	fmt.Println("Teardown TT")
}
