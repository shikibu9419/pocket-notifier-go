package api

import (
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"github.com/mattn/go-jsonpointer"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var baseUrl = "https://getpocket.com/v3/"

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

type pocketApi struct {
	env    Env
	params url.Values
}

func NewPocket() *pocketApi {
	var env Env
	envconfig.Process("pocket", &env)
	params := url.Values{"access_token": {env.AccessToken}, "consumer_key": {env.ConsumerKey}}

	return &pocketApi{
		env:    env,
		params: params}
}

func (api pocketApi) GetArticles(tag string, sort string) []Article {
	req, err := http.NewRequest("POST", baseUrl+"get", nil)
	if err != nil {
		log.Fatal(err)
	}

	api.params.Add("tag", tag)
	api.params.Add("sort", sort)
	api.params.Add("contentType", "article")
	req.URL.RawQuery = api.params.Encode()

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
		imageUrl := api.env.NoImageUrl
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

func (api pocketApi) GetRandomArticles() ([]Article, string) {
	allArticles := api.GetArticles(api.env.Tag, api.env.Sort)

	articles := make([]Article, 0, api.env.MaxCount)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < api.env.MaxCount; i++ {
		index := rand.Intn(api.env.MaxCount)
		articles = append(articles, allArticles[index])
		allArticles = append(allArticles[:index], allArticles[index+1:]...)
	}

	return articles, api.env.Tag
}
