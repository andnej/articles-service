package article_controller

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"example.com/articles/article_service"
	"github.com/gogearbox/gearbox"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var (
	gb   gearbox.Gearbox
	as   article_service.ArticleService
	port int
	json jsoniter.API = jsoniter.ConfigCompatibleWithStandardLibrary
)

func SetupTest() {
	port = rand.Intn(1000) + 5000
	gb = gearbox.New()
	as = article_service.Configure(true)
	gb.Group("/api", Setup(gb, as))
	go gb.Start(fmt.Sprintf(":%v", port))
	<-time.After(20 * time.Millisecond)
}

func CleanUpTest() {
	articles, _ := as.FindAll()
	for _, a := range articles {
		as.Delete(a.Id)
	}
	go gb.Stop()
}

func TestWhenPostThenGetHasOne(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	RetrieveAndAssertArticles(t, 0)
	createdArticle := PostAndAssertRandomArticle(t)
	if createdArticle.Id != 1 {
		t.Errorf("Expecting object created with id 1, got %v", createdArticle.Id)
	}
	RetrieveAndAssertArticles(t, 1)
}

func TestPostThreeTimesThenRetrievedHasThree(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	PostAndAssertRandomArticle(t)
	PostAndAssertRandomArticle(t)
	PostAndAssertRandomArticle(t)

	RetrieveAndAssertArticles(t, 3)
}

func TestWhenPostAndPutThenValueUpdated(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	article := PostAndAssertRandomArticle(t)
	newTitle := "Little Stuart"
	article.Title = newTitle
	updated := PutAndAssertArticle(t, article)
	if updated.Title != newTitle {
		t.Errorf("Update returned unchanged value")
	}
	RetrieveAndAssertArticles(t, 1)
}

func TestPutNonExistantIdThenStatusNotFound(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	article := randomArticle()
	article.Id = 1

	status, _, err := put(article)
	assertNotError(t, err)
	assertStatus(t, status, gearbox.StatusNotFound)
}

func TestPutInvalidIdThenStatusBadRequest(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	status, _, err := customHttpRequest("PUT", urlWithInvalidId(), randomArticle())
	assertNotError(t, err)
	assertStatus(t, status, gearbox.StatusBadRequest)
}

func TestPostThenGetHasTheSameResult(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	article := PostAndAssertRandomArticle(t)
	retrieved := GetAndAssertArticle(t, article.Id)
	if retrieved.Title != article.Title || retrieved.Body != retrieved.Body {
		t.Errorf("Get value different with Posted")
	}
}

func TestPostInvalidBodyThenGetStatusBadRequest(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	buf := []byte("{'id': '10'}")

	req := fasthttp.AcquireRequest()
	req.SetBody(buf)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url())
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	err := fasthttp.Do(req, res)
	assertNotError(t, err)
	assertStatus(t, res.StatusCode(), gearbox.StatusBadRequest)
}

func TestGetNonExistantId(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	var dst []byte
	status, body, err := fasthttp.Get(dst, urlWithId(1))
	fmt.Printf("Received: %v %v\n", status, body)
	assertNotError(t, err)
	assertStatus(t, status, gearbox.StatusNotFound)
	status, body, err = fasthttp.Get(dst, urlWithInvalidId())
	fmt.Printf("Received: %v %v\n", status, body)
	assertNotError(t, err)
	assertStatus(t, status, gearbox.StatusBadRequest)
}

func TestPostDeleteThenRetrievedHasZero(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	article := PostAndAssertRandomArticle(t)
	DeleteAndAssertArticle(t, article)
	RetrieveAndAssertArticles(t, 0)
}

func TestDeleteNonExistantIdThenGetStatusNotFound(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	status, _, err := delete(100)
	assertNotError(t, err)
	assertStatus(t, status, gearbox.StatusNotFound)
}

func TestDeleteNonExistantIdThenGetStatusBadRequest(t *testing.T) {
	SetupTest()
	defer CleanUpTest()

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("DELETE")
	req.SetRequestURI(urlWithInvalidId())
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	err := fasthttp.Do(req, res)
	assertNotError(t, err)
	assertStatus(t, res.StatusCode(), gearbox.StatusBadRequest)
}

func assertNotError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assertStatus(t *testing.T, status int, expected_status int) {
	if status != expected_status {
		t.Errorf("Status %v, Expected %v", status, expected_status)
	}
}

func assertArticles(t *testing.T, articles []*article_service.Article, expected_size int) {
	if size := len(articles); size != expected_size {
		t.Errorf("Size is %v should be %v", size, expected_size)
	}
}

func randomArticle() *article_service.Article {
	titles := []string{"Mashmallow", "Mushroom", "Miarso", "Migoyeng"}
	article := new(article_service.Article)
	article.Title = titles[rand.Intn(len(titles))]
	article.Body = fmt.Sprintf("%v", rand.Int())
	return article
}

func extractArticles(body []byte) ([]*article_service.Article, error) {
	var articles []*article_service.Article = make([]*article_service.Article, 0)
	err := json.Unmarshal(body, &articles)
	return articles, err
}

func RetrieveAndAssertArticles(t *testing.T, expected_size int) []*article_service.Article {
	var dst []byte
	status, body, err := fasthttp.Get(dst, url())
	fmt.Printf("Received: %v %v\n", status, body)
	assertNotError(t, err)
	if status != gearbox.StatusOK {
		t.Errorf("Unexpected status %v", status)
		return nil
	} else {
		articles, errx := extractArticles(body)
		assertNotError(t, errx)
		assertArticles(t, articles, expected_size)
		return articles
	}
}

func PostAndAssertRandomArticle(t *testing.T) *article_service.Article {
	article := randomArticle()
	postStatus, postBody, errp := post(article)
	fmt.Printf("Received %v %v\n", postStatus, postBody)
	assertNotError(t, errp)
	assertStatus(t, postStatus, gearbox.StatusCreated)
	createdArticle := new(article_service.Article)
	json.Unmarshal(postBody, &createdArticle)
	return createdArticle
}

func PutAndAssertArticle(t *testing.T, a *article_service.Article) *article_service.Article {
	putStatus, putBody, errp := put(a)
	fmt.Printf("Received %v %v\n", putStatus, putBody)
	assertNotError(t, errp)
	assertStatus(t, putStatus, gearbox.StatusOK)
	retrievedArticle := new(article_service.Article)
	json.Unmarshal(putBody, &retrievedArticle)
	return retrievedArticle
}

func GetAndAssertArticle(t *testing.T, id int) *article_service.Article {
	var dst []byte
	status, body, err := fasthttp.Get(dst, urlWithId(id))
	fmt.Printf("Received: %v %v\n", status, body)
	assertNotError(t, err)
	assertStatus(t, status, gearbox.StatusOK)
	retrievedArticle := new(article_service.Article)
	json.Unmarshal(body, &retrievedArticle)
	return retrievedArticle
}

func DeleteAndAssertArticle(t *testing.T, a *article_service.Article) {
	deleteStatus, deleteBody, err := delete(a.Id)
	fmt.Printf("Received %v %v\n", deleteStatus, deleteBody)
	assertNotError(t, err)
	assertStatus(t, deleteStatus, gearbox.StatusOK)
}

func post(a *article_service.Article) (status int, body []byte, err error) {
	return customHttpRequest("POST", url(), a)
}

func put(a *article_service.Article) (status int, body []byte, err error) {
	return customHttpRequest("PUT", urlWithId(a.Id), a)
}

func delete(id int) (status int, body []byte, err error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("DELETE")
	req.SetRequestURI(urlWithId(id))
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.Do(req, res); err != nil {
		return 0, nil, err
	}
	return res.StatusCode(), res.Body(), nil
}

func customHttpRequest(method string, url string, obj interface{}) (status int, body []byte, err error) {
	str, _ := json.Marshal(obj)
	buf := []byte(str)

	req := fasthttp.AcquireRequest()
	req.SetBody(buf)
	req.Header.SetMethod(method)
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.Do(req, res); err != nil {
		return 0, nil, err
	}
	return res.StatusCode(), res.Body(), nil
}

func url() string {
	return fmt.Sprintf("http://127.0.0.1:%v/api/article/", port)
}

func urlWithId(id int) string {
	return fmt.Sprintf("http://127.0.0.1:%v/api/article/%v/", port, id)
}

func urlWithInvalidId() string {
	return fmt.Sprintf("http://127.0.0.1:%v/api/article/nonexistant/", port)
}
