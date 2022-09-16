package service

import "testing"

type Vector[T any] []T
type NumSlice[T int | float64] []T
type DataMap[K string, V any] map[K]V
type Ch[T any] chan T
type MyStruct[S int | string, P map[S]string] struct {
	Name    string
	Content S
	Job     P
}

func TestIndex(t *testing.T) {
	m := make(map[int]int, 3)
	x := len(m)
	m[1] = m[1]
	y := len(m)
	t.Log("length:", x, y)
}

func TestGeneric(t *testing.T) {
	vector := Vector[string]{"a", "b", "c"}
	t.Logf("Vector: 类型=%T, val=%+v", vector, vector)

	ns := NumSlice[int]{1, 2, 3, 4, 5}
	t.Logf("NumSlice: 类型=%T, val=%+v", ns, ns)

	dm := DataMap[string, int]{
		"zx": 123,
		"as": 456,
		"qw": 789,
	}
	t.Logf("DataMap: 类型=%T, val=%+v", dm, dm)

	ch := make(Ch[int], 1)
	ch <- 10
	num := <-ch
	t.Logf("Ch: 类型=%T, val=%+v", num, num)
	t.Logf("Ch: 类型=%T, val=%+v", ch, ch)

}
