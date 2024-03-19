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
	"github.com/tmc/langchaingo/schema"
)

const PROMPT = "You are a helpful assistant describing images for blind screen reader users. Please describe this image."

var llm llms.Model

func init() {
	apiRoot := os.Getenv("OLLAMA_API")
	if apiRoot == "" {
		log.Fatal("OLLAMA_API environment variable is not defined")
		os.Exit(1)
	}
	var err error
	llm, err = ollama.New(ollama.WithServerURL(apiRoot), ollama.WithModel("llava:34b"))
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
		response, err := llm.GenerateContent(context.Background(), content)
		if err != nil {
			return err
		}
		followups, err := app.Dao().FindCollectionByNameOrId("followups")
		if err != nil {
			return err
		}
		followupIds := record.GetStringSlice("followups")
		for _, choice := range response.Choices {
			followup := models.NewRecord(followups)
			followup.Set("user", false)
			followup.Set("text", choice.Content)
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
}
