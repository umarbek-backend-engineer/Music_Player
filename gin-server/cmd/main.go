package main

import (
	cgf "gin-server/internal/config"
	"gin-server/internal/router"
)

func main() {
	port := cgf.Load().Api_Port

	r := router.Route()

	r.Run(":" + port)
}
