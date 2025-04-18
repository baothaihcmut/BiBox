package services

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
)

type IndexPayload struct {
	FileId  string `json:"file_id"`
	Content string `json:"content"`
}

type IndexTextService interface {
	IndexFile(context.Context, string, string) error
}

type ESIndexTextService struct {
	esClient *elasticsearch.Client
}

func (e *ESIndexTextService) IndexFile(ctx context.Context, fileId string, content string) error {
	doc := IndexPayload{
		FileId:  fileId,
		Content: content,
	}
	docJson, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	_, err = e.esClient.Index(
		"file_texts",
		bytes.NewReader(docJson),
		e.esClient.Index.WithDocumentID(fileId),
		e.esClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	return nil
}
