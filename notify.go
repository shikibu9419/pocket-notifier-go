package main

import (
	mySlack "./slack"
	"fmt"
	"github.com/slack-go/slack"
)

func appendArticleSections(blocks []slack.Block, itemId string) []slack.Block {
	divider := slack.NewDividerBlock()

	articleText := slack.NewTextBlockObject("mrkdwn", "総文字数: ", false, false)
	articleImage := slack.NewImageBlockElement("http://design-ec.com/d/e_others_50/l_e_others_501.png", "no_image")
	article := slack.NewSectionBlock(articleText, nil, slack.NewAccessory(articleImage))

	readButton := slack.NewButtonBlockElement("read_article", "item_id", slack.NewTextBlockObject("plain_text", "読んだ", false, false))
	read := slack.NewActionBlock("block_id", readButton)

	blocks = append(blocks, divider)
	blocks = append(blocks, article)
	blocks = append(blocks, read)
	return blocks
}

func notify() {
	header := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "この記事を読むのです...", false, false), nil, nil)
	blocks := []slack.Block{header}

	blocks = appendArticleSections(blocks, "itemId")

	msg := slack.WebhookMessage{Text: "debug message", Blocks: &slack.Blocks{BlockSet: blocks}}

	w := mySlack.NewWebhook("pocket", msg)
	err := w.Send()

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	notify()
}
