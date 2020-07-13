package pocket

import (
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"github.com/mattn/go-jsonpointer"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var baseUrl = "https://getpocket.com/v3/"
var env Env

type Env struct {
	AccessToken string `split_words:"true"`
	ConsumerKey string `split_words:"true"`
	NoImageUrl  string `split_words:"true"`
	Tag         string `default:"todo"`
	Sort        string `default:"newest"`
	MaxCount    int    `split_words:"true" default:"2"`
}

type Article struct {
	ItemId        string
	ResolvedTitle string
	ResolvedUrl   string
	ImageUrl      string
	WordCount     int
}

func GetArticles(tag string, sort string) []Article {
	envconfig.Process("pocket", &env)

	req, err := http.NewRequest("POST", baseUrl+"get", nil)
	if err != nil {
		log.Fatal(err)
	}

	params := req.URL.Query()
	params.Add("access_token", env.AccessToken)
	params.Add("consumer_key", env.ConsumerKey)
	params.Add("tag", tag)
	params.Add("sort", sort)
	params.Add("contentType", "article")

	req.URL.RawQuery = params.Encode()

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var payload interface{}
	byteArray, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(byteArray, &payload)

	list, err := jsonpointer.Get(payload, "/list")
	if err != nil {
		log.Fatal(err)
	}

	articleList := make([]interface{}, 0)
	for _, article := range list.(map[string]interface{}) {
		articleList = append(articleList, article)
	}

	articles := make([]Article, 0, len(articleList))
	for _, article := range articleList {
		articleMap := article.(map[string]interface{})
		imageUrl := env.NoImageUrl
		if articleMap["top_image_url"] != nil {
			imageUrl = articleMap["top_image_url"].(string)
		}
		wordCount, _ := strconv.Atoi(articleMap["word_count"].(string))

		articles = append(articles, Article{
			ItemId:        articleMap["item_id"].(string),
			ResolvedTitle: articleMap["resolved_title"].(string),
			ResolvedUrl:   articleMap["resolved_url"].(string),
			ImageUrl:      imageUrl,
			WordCount:     wordCount,
		})
	}

	return articles
}

func GetRandomArticles() ([]Article, string) {
	envconfig.Process("pocket", &env)

	allArticles := GetArticles(env.Tag, env.Sort)

	articles := make([]Article, 0, env.MaxCount)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < env.MaxCount; i++ {
		articles = append(articles, allArticles[rand.Intn(env.MaxCount)])
	}

	return articles, env.Tag
}
