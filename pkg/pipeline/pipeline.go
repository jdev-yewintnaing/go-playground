package pipeline

import (
	"context"
	"sync"
)

type Result[T any] struct {
	Data T
	Err  error
}

type Processor[Input any, Output any] interface {
	Process(Input) (Output, error)
}

func RunEngine[In any, Out any](
	ctx context.Context,
	items []In,
	p Processor[In, Out],
	concurrency int) []Result[Out] {

	outChannel := make(chan Result[Out], len(items))
	inChannel := make(chan In, len(items))

	var wg sync.WaitGroup
	for range concurrency {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for item := range inChannel {
				select {
				case <-ctx.Done():
					return
				default:
					out, err := p.Process(item)
					outChannel <- Result[Out]{
						Data: out,
						Err:  err,
					}
				}
			}
		}()
	}

	// for _, item := range items {
	// 	inChannel <- item
	// }

	go func() {
		defer close(inChannel)
		for _, item := range items {
			select {
			case <-ctx.Done():
				return
			case inChannel <- item:
			}
		}
	}()

	wg.Wait()

	close(outChannel)

	var results []Result[Out]
	for result := range outChannel {
		results = append(results, result)
	}

	return results

}
