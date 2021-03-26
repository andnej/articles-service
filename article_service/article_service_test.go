package article_service

import (
	"testing"
)

func TestInit(t *testing.T) {
	size := Len()
	if size != 0 {
		t.Errorf("Initialized but not zero length")
	}
}

func TestSave(t *testing.T) {
	Reset()
	article := new(Article)
	article.Title = "Holistic"
	article.Body = "The Body"
	savedArticle, err := Save(article)
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
	Reset()
	article := new(Article)
	article.Title = "Holistic"
	article.Body = "The Body"
	firstArticle, err1 := Save(article)
	secondArticle, err2 := Save(article)
	if err1 != nil || err2 != nil {
		t.Errorf("Error occurred")
	}
	if firstArticle.Id != 1 {
		t.Errorf("First Article id is %v, should be 1", firstArticle.Id)
	}
	if secondArticle.Id != 2 {
		t.Errorf("Second Article id is %v, should be 2", secondArticle.Id)
	}
	firstArticle.Title = "Babushka"
	updatedArticle, _ := Save(firstArticle)
	if updatedArticle.Id != firstArticle.Id {
		t.Errorf("Saving duplicate item should update instead of storing another copy")
	}
}

func TestFindAll(t *testing.T) {
	Reset()
	storedArticles, err := FindAll()
	if err != nil {
		t.Errorf("Error findAll()")
	}
	if storedArticles == nil {
		t.Errorf("returned value should not be nil")
	}
	if len(storedArticles) != Len() {
		t.Errorf("different length between Len() and returned")
	}

	Save(new(Article))
	if len(storedArticles) == Len() {
		t.Errorf("Unsafe returned value")
	}
}

func TestFindOne(t *testing.T) {
	Reset()
	firstArticle, _ := Save(new(Article))

	found, err := FindOne(int(firstArticle.Id))
	if found == nil || err != nil {
		t.Errorf("Error FindOne %v", err)
	}

	found, err = FindOne(2)
	if found != nil || err == nil {
		t.Errorf("Should not find one")
	}
}

func TestDelete(t *testing.T) {
	Reset()
	Save(new(Article))
	Save(new(Article))
	Save(new(Article))

	deleted, err := Delete(2)
	if err != nil {
		t.Errorf("Error delete %v", err)
	}
	if deleted.Id != 2 {
		t.Errorf("Delete wrong item")
	}
	if Len() != 2 {
		t.Errorf("Invalid item count after deletion")
	}
}
