package main

import (
	"fmt"
	"go-playground/pkg/generic"
	"go-playground/pkg/scheduler"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	manager := scheduler.NewManager()

	manager.Add(scheduler.EmailTask{
		Email: "yewintnaing@gmail.com",
	})
	manager.Add(scheduler.EmailTask{
		Email: "yewintnaing44@gmail.com",
	})

	err := manager.RunAll()
	if err != nil {
		return
	}

	// Filtering Ints
	nums := []int{1, 2, 3, 4, 5}
	evens := generic.Filter(nums, func(n int) bool { return n%2 == 0 })

	fmt.Println(evens)

}
