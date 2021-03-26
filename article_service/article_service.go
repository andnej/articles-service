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
	articles      []*Article  = []*Article{}
	nextId        int         = 1
	mu            *sync.Mutex = &sync.Mutex{}
	NOT_FOUND     int         = -1
	ERR_NOT_FOUND error       = errors.New("Specified id not found")
)

func Reset() {
	articles = []*Article{}
	nextId = 1
}

func Save(article *Article) (*Article, error) {
	mu.Lock()
	defer mu.Unlock()
	if article.Id < 1 {
		newArticle := new(Article)
		newArticle.Id = nextId
		newArticle.Title = article.Title
		newArticle.Body = article.Body
		articles = append(articles, newArticle)
		nextId++
		return newArticle, nil
	} else if theIndex := index(article.Id); theIndex >= 0 {
		oldArticle := articles[theIndex]
		oldArticle.Title = article.Title
		oldArticle.Body = article.Body
		return oldArticle, nil
	} else {
		return nil, ERR_NOT_FOUND
	}
}

func FindAll() ([]*Article, error) {
	return articles, nil
}

func FindOne(id int) (*Article, error) {
	theIndex := index(id)
	if theIndex == NOT_FOUND {
		return nil, ERR_NOT_FOUND
	} else {
		return articles[theIndex], nil
	}
}

func Delete(id int) (*Article, error) {
	mu.Lock()
	defer mu.Unlock()
	theIndex := index(id)
	if theIndex == NOT_FOUND {
		return nil, ERR_NOT_FOUND
	} else {
		result := articles[theIndex]
		articles = append(articles[:theIndex], articles[theIndex+1:]...)
		return result, nil
	}
}

func index(id int) int {
	for index, art := range articles {
		if art.Id == id {
			return index
		}
	}

	return NOT_FOUND
}

func Len() int {
	return len(articles)
}
