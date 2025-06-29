package rabbitmq

import (
	"errors"
	"fmt"
	"sync"
)

type TestingConsumer struct {
	mutex sync.Mutex
}

func (consume *TestingConsumer) Consume(message any) error {
	consume.mutex.Lock()
	defer consume.mutex.Unlock()

	data, ok := message.(map[string]interface{})
	if !ok {
		return errors.New("Your message is not map[string]interface{}")
	}

	fmt.Println(data["name"])
	return nil
}
