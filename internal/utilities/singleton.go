package utilities

import (
	"fmt"
	"sync"
)

func Singleton[T any](init func() (T, error)) func() T {
	var (
		once sync.Once
		t    T
	)

	return func() T {
		once.Do(func() {
			var err error
			t, err = init()
			if err != nil {
				panic(fmt.Sprintf("creating %T failed: %s", t, err))
			}
		})

		return t
	}
}
