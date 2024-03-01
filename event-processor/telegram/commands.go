package telegram

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"school-project/data"
	"school-project/summary"
	"strings"
)

const (
	RndCmd   = "/random"
	HelpCmd  = "/help"
	StartCmd = "/start"
	LastCmd  = "/last"
	ClearCmd = "/clear"
	FirstCmd = "/first"
)

func (p *Processor) doCmd(text string, chatID int, username string, summarizer *summary.OpenAISummarizer, prompt string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username, summarizer, prompt)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case LastCmd:
		return p.sendLast(chatID, username, summarizer, prompt)
	case ClearCmd:
		return p.clearAll(chatID, username)
	case FirstCmd:
		return p.sendFirst(chatID, username, summarizer, prompt)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {

	page := &data.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.data.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.data.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string, summarizer *summary.OpenAISummarizer, prompt string) (err error) {

	log.Print("article is creating")
	p.tg.SendMessage(chatID, "Секундочку...")

	page, err := p.data.PickRandom(username)
	if err != nil && !errors.Is(err, data.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, data.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	article, err := summarizer.Summarize(prompt + page.URL)
	log.Print(article)
	if err != nil {
		return err
	}
	err = p.tg.SendMessage(chatID, article+"\n"+page.URL)
	if err != nil {
		return err
	}
	log.Print("article sent")
	return p.data.Remove(page)
}

func (p *Processor) clearAll(chatID int, username string) (err error) {

	err = p.tg.SendMessage(chatID, "Секундочку...")
	if err != nil {
		fmt.Println(err)
	}

	err = p.data.ClearAll(username)
	p.tg.SendMessage(chatID, "Статьи удалены")
	return
}

func (p *Processor) sendLast(chatID int, username string, summarizer *summary.OpenAISummarizer, prompt string) (err error) {

	log.Print("article is creating")
	p.tg.SendMessage(chatID, "Секундочку...")

	page, err := p.data.PickLast(username)
	if err != nil && !errors.Is(err, data.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, data.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	article, err := summarizer.Summarize(prompt + page.URL)
	log.Print(article)
	if err != nil {
		return err
	}
	err = p.tg.SendMessage(chatID, article+"\n"+page.URL)
	if err != nil {
		return err
	}
	log.Print("article sent")
	return p.data.Remove(page)
}

func (p *Processor) sendFirst(chatID int, username string, summarizer *summary.OpenAISummarizer, prompt string) (err error) {

	log.Print("article is creating")
	p.tg.SendMessage(chatID, "Секундочку...")

	page, err := p.data.PickFirst(username)
	if err != nil && !errors.Is(err, data.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, data.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	article, err := summarizer.Summarize(prompt + page.URL)
	log.Print(article)
	if err != nil {
		return err
	}
	err = p.tg.SendMessage(chatID, article+"\n"+page.URL)
	if err != nil {
		return err
	}
	log.Print("article sent")
	return p.data.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp+msgSendMe)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

const msgHelp = "Этот бот может сохранять твои ссылки на статьи и по запросу отправлять тебе ссылку и краткую выжимку на статью по данной ссылке, сгенерированной искусственным интеллектом (gpt 3.5 turbo)."

const msgHello = "Привет! \n\n" + msgHelp + msgSendMe

const (
	msgUnknownCommand = "Неизвестная команда"
	msgNoSavedPages   = "Сохраненных статей не найдено"
	msgSaved          = "Сохранено!"
	msgAlreadyExists  = "Эта статья уже существует в Вашем списке"
	msgSendMe         = "\n\nДля начала просто отправь мне ссылку на статью, она автоматически сохранится.\n\nДалее просто используй команды из списка."
)
