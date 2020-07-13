package pocket

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/mattn/go-jsonpointer"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var baseUrl = "https://getpocket.com/v3/"
var env Env

type Env struct {
	AccessToken string `split_words:"true"`
	ConsumerKey string `split_words:"true"`
	NoImageUrl  string `split_words:"true"`
}

type Article struct {
	ItemId        string
	ResolvedTitle string
	ResolvedUrl   string
	ImageUrl      string
	WordCount     int
}

func GetArticles(tag string, sort string, count int) []Article {
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

	if sort == "random" {
		params.Add("sort", "newest")
	}
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
	for _, value := range list.(map[string]interface{}) {
		articleList = append(articleList, value)
	}

	articles := make([]Article, 0, count)
	for _, article := range articleList[:count] {
		articleMap := article.(map[string]interface{})
		imageUrl := articleMap["top_image_url"].(string)
		if imageUrl == "" {
			imageUrl = env.NoImageUrl
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
