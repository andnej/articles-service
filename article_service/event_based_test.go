package article_service

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestConfigureEventBus(t *testing.T) {
	eventBus, err := ConfigureEventBus(true)
	if err != nil {
		t.Error(err)
	}
	resultChannel := make(chan FindAllResponse)
	defer close(resultChannel)
	req := FindAllRequest{
		resultChannel: resultChannel,
	}
	eventBus.FindAllRequestChannel() <- req
	var findAllResponse *FindAllResponse
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("Spawning request")
	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(500 * time.Millisecond):
				fmt.Println("timeout")
				return
			case res := <-resultChannel:
				findAllResponse = &res
				return
			}
		}
	}()
	wg.Wait()
	if findAllResponse == nil {
		t.Errorf("No response in 500 ms")
	}
}
