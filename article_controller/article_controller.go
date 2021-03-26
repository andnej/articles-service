package article_controller

import (
	"strconv"

	"example.com/articles/article_service"
	"github.com/gogearbox/gearbox"
)

const ARTICLE_GROUP_NAME string = "/article"

func Setup(gb gearbox.Gearbox, as article_service.ArticleService) []*gearbox.Route {
	extractId := func(ctx gearbox.Context) (int, error) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return -1, err
		}
		return id, nil
	}

	routes := []*gearbox.Route{
		gb.Get("/", func(ctx gearbox.Context) {
			articles, err := as.FindAll()
			if err != nil {
				ctx.Status(gearbox.StatusInternalServerError)
			} else {
				ctx.SendJSON(articles)
			}
		}),
		gb.Get("/:id", func(ctx gearbox.Context) {
			id, err := extractId(ctx)
			if err != nil {
				ctx.Status(gearbox.StatusBadRequest)
				return
			}
			article, errfo := as.FindOne(id)
			if errfo != nil {
				ctx.Status(gearbox.StatusNotFound)
				return
			}
			ctx.SendJSON(article)
		}),
		gb.Put("/:id", func(ctx gearbox.Context) {
			id, err := extractId(ctx)
			if err != nil {
				ctx.Status(gearbox.StatusBadRequest)
				return
			}
			oldArticle, errfo := as.FindOne(id)
			if errfo != nil {
				ctx.Status(gearbox.StatusNotFound)
				return
			}
			newArticle := new(article_service.Article)
			err = ctx.ParseBody(newArticle)
			oldArticle.Title = newArticle.Title
			oldArticle.Body = newArticle.Body
			as.Save(oldArticle)
			ctx.SendJSON(oldArticle)
		}),
		gb.Post("/", func(ctx gearbox.Context) {
			article := new(article_service.Article)
			err := ctx.ParseBody(article)
			if err != nil {
				ctx.Status(gearbox.StatusBadRequest)
				return
			}
			article.Id = 0
			if created, err := as.Save(article); err == nil {
				ctx.Status(gearbox.StatusCreated)
				ctx.SendJSON(created)
			} else {
				ctx.Status(gearbox.StatusInternalServerError)
			}
		}),
		gb.Delete("/:id", func(ctx gearbox.Context) {
			id, err := extractId(ctx)
			if err != nil {
				ctx.Status(gearbox.StatusBadRequest)
				return
			}
			deleted, errd := as.Delete(id)
			if errd != nil {
				ctx.Status(gearbox.StatusNotFound)
			} else {
				ctx.SendJSON(deleted)
			}
		}),
	}

	articlesRoutes := gb.Group(ARTICLE_GROUP_NAME, routes)
	return articlesRoutes
}
