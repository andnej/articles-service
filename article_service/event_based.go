package article_service

import (
	"github.com/ztrue/shutdown"
	"gopkg.in/thejerf/suture.v4"
)

type FindAllRequest struct {
	resultChannel chan<- FindAllResponse
}

type FindAllResponse struct {
	articles []*Article
	err      error
}

type ArticleEventBus struct {
	findAllRequestChannel chan FindAllRequest
	articleService        ArticleService
}

func (eb *ArticleEventBus) Serve() {
	var req FindAllRequest
	for {
		select {
		case req = <-eb.findAllRequestChannel:
			articles, err := eb.articleService.FindAll()
			findAllResponse := FindAllResponse{
				articles: articles,
				err:      err,
			}
			req.resultChannel <- findAllResponse
		}
	}
}

func (eb *ArticleEventBus) Stop() {

}

func (eb *ArticleEventBus) FindAllRequestChannel() chan<- FindAllRequest {
	return eb.findAllRequestChannel
}

type ArticleEventBusClient interface {
	FindAllRequestChannel() chan<- FindAllRequest
}

func ConfigureEventBus(test bool) (ArticleEventBusClient, error) {
	supervisor := suture.NewSimple("ArticleEventBusSupervisor")
	eb := ArticleEventBus{
		findAllRequestChannel: make(chan FindAllRequest, 1),
		articleService:        Configure(test),
	}

	service := &eb
	supervisor.Add(service)
	go supervisor.ServeBackground()
	shutdown.Add(func() {
		supervisor.Stop()
	})
	var ebc ArticleEventBusClient
	ebc = &eb

	return ebc, nil
}
