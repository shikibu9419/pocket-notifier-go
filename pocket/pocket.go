package pocket

import (
	"encoding/json"
	"fmt"
	"github.com/mattn/go-jsonpointer"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"os"
)

var baseUrl = "https://getpocket.com/v3/"

type Article struct {
	ItemId        string
	ResolvedTitle string
	ResolvedUrl   string
	WordCount     int
}

func GetArticles(tag string, sort string, count int) []Article {
	req, err := http.NewRequest("POST", baseUrl+"get", nil)
	if err != nil {
		log.Fatal(err)
	}

	params := req.URL.Query()
	params.Add("access_token", os.Getenv("POCKET_ACCESS_TOKEN"))
	params.Add("consumer_key", os.Getenv("POCKET_CONSUMER_KEY"))
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
		wordCount, _ := strconv.Atoi(articleMap["word_count"].(string))

		articles = append(articles, Article{
			ItemId:        articleMap["item_id"].(string),
			ResolvedTitle: articleMap["resolved_title"].(string),
			ResolvedUrl:   articleMap["resolved_url"].(string),
			WordCount:     wordCount,
		})
	}

	return articles
}

func main() {
	fmt.Println(GetArticles("todo", "newest", 2))
}
