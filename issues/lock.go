package issues

import (
	"fmt"
	"sync"
)

func mutexParam() {
	mutex := sync.Mutex{}
	mutex.Lock()
	
	handler := func(m *sync.Mutex) {
		m.Unlock()
		fmt.Println("handler")
	}
	
	handler(&mutex)
	// go handler(&mutex)
	// mutex.Unlock()
	fmt.Println("main")
}