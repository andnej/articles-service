package article_service

import (
	"errors"
	"sync"
)

type Article struct {
	Id    int
	Title string
	Body  string
}

var (
	NOT_FOUND     int   = -1
	ERR_NOT_FOUND error = errors.New("Specified id not found")
)

type ArticleService interface {
	Save(article *Article) (*Article, error)
	FindAll() ([]*Article, error)
	FindOne(id int) (*Article, error)
	Delete(id int) (*Article, error)
	Len() int
}

type InMemoryArticleService struct {
	articles []*Article
	nextId   int
	mu       *sync.Mutex
}

func Configure(test bool) ArticleService {
	var as ArticleService
	as = &InMemoryArticleService{
		articles: []*Article{},
		nextId:   1,
		mu:       &sync.Mutex{},
	}

	return as
}

func (im *InMemoryArticleService) Save(article *Article) (*Article, error) {
	im.mu.Lock()
	defer im.mu.Unlock()
	if article.Id < 1 {
		newArticle := new(Article)
		newArticle.Id = im.nextId
		newArticle.Title = article.Title
		newArticle.Body = article.Body
		im.articles = append(im.articles, newArticle)
		im.nextId++
		return newArticle, nil
	} else if theIndex := im.index(article.Id); theIndex >= 0 {
		oldArticle := im.articles[theIndex]
		oldArticle.Title = article.Title
		oldArticle.Body = article.Body
		return oldArticle, nil
	} else {
		return nil, ERR_NOT_FOUND
	}
}

func (im *InMemoryArticleService) FindAll() ([]*Article, error) {
	return im.articles, nil
}

func (im *InMemoryArticleService) FindOne(id int) (*Article, error) {
	theIndex := im.index(id)
	if theIndex == NOT_FOUND {
		return nil, ERR_NOT_FOUND
	} else {
		return im.articles[theIndex], nil
	}
}

func (im *InMemoryArticleService) Delete(id int) (*Article, error) {
	im.mu.Lock()
	defer im.mu.Unlock()
	theIndex := im.index(id)
	if theIndex == NOT_FOUND {
		return nil, ERR_NOT_FOUND
	} else {
		result := im.articles[theIndex]
		im.articles = append(im.articles[:theIndex], im.articles[theIndex+1:]...)
		return result, nil
	}
}

func (im *InMemoryArticleService) index(id int) int {
	for index, art := range im.articles {
		if art.Id == id {
			return index
		}
	}

	return NOT_FOUND
}

func (im *InMemoryArticleService) Len() int {
	return len(im.articles)
}
