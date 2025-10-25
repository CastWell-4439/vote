package main

import (
	"sync"
	"toupiao/router"
	"toupiao/service"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		service.Start()
	}()
	routers := router.Router()
	routers.Run(":8888")
}
