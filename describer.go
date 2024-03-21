package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

const PROMPT = "You are a helpful assistant describing images for blind screen reader users. Please describe this image."

var llm *llms.Model

func init() {
	apiRoot := os.Getenv("OLLAMA_API")
	if apiRoot == "" {
		log.Println("OLLAMA_API environment variable is not defined, using stub data instead")
		return
	}
	var err error
	*llm, err = ollama.New(ollama.WithServerURL(apiRoot), ollama.WithModel("llava:34b"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func initializeDescriber(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateRequest("images").Add(func(e *core.RecordCreateEvent) error {
		record := e.Record
		info := apis.RequestInfo(e.HttpContext)
		user := info.AuthRecord
		if user == nil {
			return errors.New("Unauthorized")
		}
		images := user.GetStringSlice("images")
		images = append(images, record.Id)
		user.Set("images", images)
		err := app.Dao().SaveRecord(user)
		if err != nil {
			return err
		}
		key := record.BaseFilesPath() + "/" + record.GetString("file")
		fsys, err := app.NewFilesystem()
		if err != nil {
			return err
		}
		defer fsys.Close()
		blob, err := fsys.GetFile(key)
		if err != nil {
			return err
		}
		defer blob.Close()
		bytes, err := io.ReadAll(blob)
		if err != nil {
			return err
		}
		content := []llms.MessageContent{
			llms.TextParts(schema.ChatMessageTypeSystem, PROMPT),
			{
				Role:  schema.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{llms.BinaryPart(blob.ContentType(), bytes)},
			},
		}
		followups, err := app.Dao().FindCollectionByNameOrId("followups")
		followupIds := record.GetStringSlice("followups")
		if llm != nil {
			response, err := (*llm).GenerateContent(context.Background(), content)
			if err != nil {
				return err
			}
			for _, choice := range response.Choices {
				followup := models.NewRecord(followups)
				followup.Set("user", false)
				followup.Set("text", choice.Content)
				if err := app.Dao().SaveRecord(followup); err != nil {
					return err
				}
				followupIds = append(followupIds, followup.Id)
			}
		} else {
			followup := models.NewRecord(followups)
			followup.Set("user", false)
			followup.Set("text", "This is a stub initial followup because `OLLAMA_API` is not set.")
			if err := app.Dao().SaveRecord(followup); err != nil {
				return err
			}
			followupIds = append(followupIds, followup.Id)
		}
		record.Set("followups", followupIds)
		err = app.Dao().SaveRecord(record)
		if err != nil {
			return err
		}
		return nil
	})
	app.OnRecordBeforeCreateRequest("followups").Add(func(e *core.RecordCreateEvent) error {
		e.Record.Set("user", true)
		return nil
	})
	app.OnRecordAfterUpdateRequest("images").Add(func(e *core.RecordUpdateEvent) error {
		record := e.Record
		ctx := context.Background()
		mem := memory.NewConversationBuffer()
		mem.ChatHistory.AddMessage(ctx, schema.SystemChatMessage{Content: PROMPT})
		followupIds := record.GetStringSlice("followups")
		for _, followupId := range followupIds {
			followup, err := app.Dao().FindRecordById("followups", followupId)
			if err != nil {
				return err
			}
			text := followup.GetString("text")
			if followup.GetBool("user") {
				mem.ChatHistory.AddUserMessage(ctx, text)
			} else {
				mem.ChatHistory.AddAIMessage(ctx, text)
			}
		}
		// chain := chains.NewConversation(llm, mem)
		// chain.Run
		return nil
	})
}
