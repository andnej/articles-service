package main

import (
	"example.com/articles/article_controller"
	"github.com/gogearbox/gearbox"
)

func main() {
	gb := gearbox.New()

	gb.Group("/api", article_controller.Setup(gb))

	gb.Start(":3000")
}
