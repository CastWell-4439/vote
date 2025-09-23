package main

import (
	"toupiao/router"
)

func main() {
	router := router.Router()
	router.Run(":8888")
}
