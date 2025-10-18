package main

import (
	"toupiao/router"
)

func main() {
	routers := router.Router()
	routers.Run(":8888")
}
