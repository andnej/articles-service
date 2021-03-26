package main

import (
	"example.com/articles/article_controller"
	"example.com/articles/article_service"
	"github.com/gogearbox/gearbox"
)

func main() {
	gb := gearbox.New()

	gb.Group("/api", article_controller.Setup(gb, article_service.Configure(true)))

	gb.Start(":3000")
}
