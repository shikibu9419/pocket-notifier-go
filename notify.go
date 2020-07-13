package main

import (
	"./api"
	"fmt"
	"github.com/slack-go/slack"
	"log"
)

func appendArticleSections(blocks []slack.Block, article api.Article) []slack.Block {
	blocks = append(blocks, slack.NewDividerBlock())

	articleText := fmt.Sprintf("%s\n%s\n総文字数: %d", article.ResolvedTitle, article.ResolvedUrl, article.WordCount)
	articleTextBlock := slack.NewTextBlockObject("mrkdwn", articleText, false, false)
	articleImage := slack.NewImageBlockElement(article.ImageUrl, article.ResolvedTitle[:15])

	blocks = append(blocks, slack.NewSectionBlock(articleTextBlock, nil, slack.NewAccessory(articleImage)))

	readButton := slack.NewButtonBlockElement("read_article", article.ItemId, slack.NewTextBlockObject("plain_text", "読んだ", false, false))

	blocks = append(blocks, slack.NewActionBlock(article.ItemId, readButton))

	return blocks
}

func notify() {
	p := api.NewPocket()
	articles, tag := p.GetRandomArticles()
	fmt.Println(articles)

	headerText := fmt.Sprintf("この記事を読むのです...\nタグ: *%s*", tag)
	header := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", headerText, false, false), nil, nil)
	blocks := []slack.Block{header}

	for _, article := range articles {
		blocks = appendArticleSections(blocks, article)
	}

	msg := slack.WebhookMessage{Text: "本日のお告げが届きました", Blocks: &slack.Blocks{BlockSet: blocks}}

	w := api.NewSlackWebhook("pocket", msg)
	err := w.Send()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	notify()
}
