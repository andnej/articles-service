package article_service

import (
	"testing"
)

var (
	articleService ArticleService
)

func setup() {
	articleService = Configure(true)
}

func TestInit(t *testing.T) {
	setup()

	size := articleService.Len()
	if size != 0 {
		t.Errorf("Initialized but not zero length")
	}
}

func TestSave(t *testing.T) {
	setup()

	article := new(Article)
	article.Title = "Holistic"
	article.Body = "The Body"
	savedArticle, err := articleService.Save(article)
	if err != nil {
		t.Errorf("Error saving article")
	}
	if savedArticle.Id == 0 {
		t.Errorf("Fail generating id")
	}
	if savedArticle.Title != article.Title || savedArticle.Body != article.Body {
		t.Errorf("different article returned")
	}
}

func TestDoubleSave(t *testing.T) {
	setup()

	firstId := articleService.Len() + 1

	article := new(Article)
	article.Title = "Holistic"
	article.Body = "The Body"
	firstArticle, err1 := articleService.Save(article)
	secondArticle, err2 := articleService.Save(article)
	if err1 != nil || err2 != nil {
		t.Errorf("Error occurred")
	}
	if firstArticle.Id != firstId {
		t.Errorf("First Article id is %v, should be %v", firstArticle.Id, firstId)
	}
	secondId := firstId + 1
	if secondArticle.Id != secondId {
		t.Errorf("Second Article id is %v, should be %v", secondArticle.Id, secondId)
	}
	firstArticle.Title = "Babushka"
	updatedArticle, _ := articleService.Save(firstArticle)
	if updatedArticle.Id != firstArticle.Id {
		t.Errorf("Saving duplicate item should update instead of storing another copy")
	}
}

func TestFindAll(t *testing.T) {
	setup()

	storedArticles, err := articleService.FindAll()
	if err != nil {
		t.Errorf("Error findAll()")
	}
	if storedArticles == nil {
		t.Errorf("returned value should not be nil")
	}
	if len(storedArticles) != articleService.Len() {
		t.Errorf("different length between Len() and returned")
	}

	articleService.Save(new(Article))
	if len(storedArticles) == articleService.Len() {
		t.Errorf("Unsafe returned value")
	}
}

func TestFindOne(t *testing.T) {
	setup()

	firstArticle, _ := articleService.Save(new(Article))

	found, err := articleService.FindOne(int(firstArticle.Id))
	if found == nil || err != nil {
		t.Errorf("Error FindOne %v", err)
	}

	size := articleService.Len()

	found, err = articleService.FindOne(size + 1)
	if found != nil || err == nil {
		t.Errorf("Should not find one")
	}
}

func TestDelete(t *testing.T) {
	setup()

	articleService.Save(new(Article))
	articleService.Save(new(Article))
	articleService.Save(new(Article))

	size := articleService.Len()
	toDelete := size - 2

	deleted, err := articleService.Delete(toDelete)
	if err != nil {
		t.Errorf("Error delete %v", err)
	}
	if deleted.Id != toDelete {
		t.Errorf("Delete wrong item")
	}
	if articleService.Len() != (size - 1) {
		t.Errorf("Invalid item count after deletion")
	}
}
