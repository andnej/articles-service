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
	if eventBus.FindAllRequestChannel() == nil {
		t.Errorf("returned event bus has nil findAllRequestChannel")
	}
	if eventBus.FindOneRequestChannel() == nil {
		t.Errorf("returned event bus has nil findOneRequestChannel")
	}
}

func TestDoubleConfigureEventBusThenResultShouldBeTheSame(t *testing.T) {
	firstEventBus, err1 := ConfigureEventBus(true)
	if err1 != nil {
		t.Error(err1)
	}
	secondEventBus, err2 := ConfigureEventBus(true)
	if err2 != nil {
		t.Error(err2)
	}
	if firstEventBus != secondEventBus {
		t.Errorf("event bus not identical")
	}
}

func TestSendingFindAllRequest(t *testing.T) {
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

func TestSendingFindOneRequest(t *testing.T) {
	eventBus, err := ConfigureEventBus(true)
	if err != nil {
		t.Error(err)
	}
	resultChannel := make(chan FindOneResponse)
	defer close(resultChannel)
	req := FindOneRequest{
		resultChannel: resultChannel,
	}
	eventBus.FindOneRequestChannel() <- req
	var findOneResponse *FindOneResponse
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
				findOneResponse = &res
				return
			}
		}
	}()
	wg.Wait()
	if findOneResponse == nil {
		t.Errorf("No response in 500 ms")
	}
}

func TestManualShutdownEventBus(t *testing.T) {
	eventBus, err := ConfigureEventBus(true)
	if err != nil {
		t.Error(err)
	}
	resultChannel := make(chan FindAllResponse)
	defer close(resultChannel)
	err = ShutdownEventBus()
	if err != nil {
		t.Error(err)
	}
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
	if findAllResponse != nil {
		t.Errorf("There is still a response even though we have shut it down")
	}
}

func TestShutdownAndReconfigure(t *testing.T) {
	firstEventBus, _ := ConfigureEventBus(true)
	ShutdownEventBus()
	secondEventBus, _ := ConfigureEventBus(true)
	if firstEventBus == secondEventBus {
		t.Errorf("second event bus is identical to before shutdown")
	}
}
