package article_service

import (
	"context"
	"time"

	"github.com/ztrue/shutdown"
	"gopkg.in/thejerf/suture.v4"
)

var supervisor *suture.Supervisor
var serviceToken suture.ServiceToken
var ebc ArticleEventBusClient

type FindAllRequest struct {
	resultChannel chan<- FindAllResponse
}

type FindAllResponse struct {
	articles []*Article
	err      error
}

type FindOneRequest struct {
	id            int
	resultChannel chan<- FindOneResponse
}

type FindOneResponse struct {
	article *Article
	err     error
}

type ArticleEventBus struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	findAllRequestChannel chan FindAllRequest
	findOneRequestChannel chan FindOneRequest
	articleService        ArticleService
}

func (eb *ArticleEventBus) Serve() {
	var far FindAllRequest
	var foner FindOneRequest
	for {
		select {
		case <-eb.ctx.Done():
			return
		case far = <-eb.findAllRequestChannel:
			articles, err := eb.articleService.FindAll()
			findAllResponse := FindAllResponse{
				articles: articles,
				err:      err,
			}
			far.resultChannel <- findAllResponse
		case foner = <-eb.findOneRequestChannel:
			article, err := eb.articleService.FindOne(foner.id)
			findOneResponse := FindOneResponse{
				article: article,
				err:     err,
			}
			foner.resultChannel <- findOneResponse
		}
	}
}

func (eb *ArticleEventBus) Stop() {
	eb.cancel()
}

func (eb *ArticleEventBus) FindAllRequestChannel() chan<- FindAllRequest {
	return eb.findAllRequestChannel
}

func (eb *ArticleEventBus) FindOneRequestChannel() chan<- FindOneRequest {
	return eb.findOneRequestChannel
}

type ArticleEventBusClient interface {
	FindAllRequestChannel() chan<- FindAllRequest
	FindOneRequestChannel() chan<- FindOneRequest
}

func ConfigureEventBus(test bool) (ArticleEventBusClient, error) {
	if supervisor == nil {
		supervisor = suture.NewSimple("ArticleEventBusSupervisor")
		go supervisor.ServeBackground()
	}
	if ebc == nil {
		ctx, cancel := context.WithCancel(context.Background())
		eb := ArticleEventBus{
			ctx:                   ctx,
			cancel:                cancel,
			findAllRequestChannel: make(chan FindAllRequest, 1),
			findOneRequestChannel: make(chan FindOneRequest, 1),
			articleService:        Configure(test),
		}

		service := &eb
		serviceToken = supervisor.Add(service)
		shutdown.Add(func() {
			supervisor.Stop()
		})
		ebc = &eb
	}

	return ebc, nil
}

func ShutdownEventBus() error {
	err := supervisor.RemoveAndWait(serviceToken, 1*time.Second)
	if err != nil {
		return err
	}
	ebc = nil
	return nil
}
