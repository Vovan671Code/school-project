package telegram

import (
	"errors"
	"fmt"
	telegram "school-project/client"
	"school-project/data"
	event_processor "school-project/event-processor"
	"school-project/summary"
)

type Processor struct {
	tg         *telegram.Client
	offset     int
	data       data.Storage
	summarizer *summary.OpenAISummarizer
	prompt     string
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, data data.Storage, summarizer *summary.OpenAISummarizer, prompt string) *Processor {
	return &Processor{
		tg:         client,
		data:       data,
		summarizer: summarizer,
		prompt:     prompt,
	}
}

func (p *Processor) Fetch(limit int) ([]event_processor.Event, error) {
	updates, err := p.tg.GetMessages(p.offset, limit)
	if err != nil {

		return nil, fmt.Errorf("can't get events: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]event_processor.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event event_processor.Event) error {
	switch event.Type {
	case event_processor.Message:
		return p.processMessage(event)
	default:
		return fmt.Errorf("can't process message: %w", ErrUnknownMetaType)
	}
}

func (p *Processor) processMessage(event event_processor.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username, p.summarizer, p.prompt); err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}
	
	return nil
}

func meta(event event_processor.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {

		return Meta{}, fmt.Errorf("can't get meta: %w", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Message) event_processor.Event {
	updType := fetchType(upd)

	res := event_processor.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == event_processor.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Message) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Message) event_processor.Type {
	if upd.Message == nil {
		return event_processor.Unknown
	}

	return event_processor.Message
}
