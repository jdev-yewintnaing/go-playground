package generic

import "testing"

type person struct {
	Id   int
	Name string
	Age  int
}

func TestCache_Get(t *testing.T) {

	p := person{
		Id:   1,
		Name: "John Doe",
		Age:  42,
	}

	var cache = NewCache[person]()

	cache.Set("1", p)

	data, _ := cache.Get("1")

	if data != p {
		t.Errorf("got %v, want %v", data, p)
	}
}

func TestCache_Concurrent(t *testing.T) {
	cache := NewCache[int]()
	const goroutines = 100
	const increments = 1000
	const expectedIncrements = increments - 1

	done := make(chan bool)

	for i := 0; i < goroutines; i++ {
		go func() {
			for j := 0; j < increments; j++ {
				cache.Set("key", j)
				_, _ = cache.Get("key")
			}
			done <- true
		}()
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}

	finalResult, _ := cache.Get("key")
	if finalResult != expectedIncrements {
		t.Errorf("got %v, want %v", finalResult, expectedIncrements)
	}

}
