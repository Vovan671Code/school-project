package main

import (
	"log"
	tgClient "school-project/client"
	"school-project/consumer/while_consumer"
	"school-project/data/files"
	"school-project/event-processor/telegram"
	"school-project/summary"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
	model       = "gpt-3.5-turbo"
	prompt      = "Сделай краткую выжимку статьи объемом 200 слов на сайте  "
	//put your ai token here:
	aiToken = ""
	//put your tg token here:
	tgToken = ""
)

// to start bot put your tokens above and run school-project.exe in terminal
func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, tgToken),
		files.New(storagePath),
		summary.NewOpenAISummarizer(aiToken, model, prompt),
		prompt,
	)

	log.Print("service started")

	consumer := while_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
