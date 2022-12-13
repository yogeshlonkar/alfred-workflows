package search

import (
	"errors"
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/rs/zerolog/log"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
)

const rwRwRw = 0666

type Service interface {
	Match(keyword string) ([]alfred.Item, error)
	Load(suggestion ...alfred.Item) error
	Close() error
	Loaded() bool
}

type SuggestionIndexMeta struct {
	UID   string `json:"uid,omitempty"`
	Data  string `json:"data,omitempty"` // original source document - json.Marshal(content)
	Match string `json:"match,omitempty"`
}

func (m *SuggestionIndexMeta) Type() string {
	return "content"
}

func (m *SuggestionIndexMeta) BleveType() string {
	return "content"
}

func NewService(indexPath string) (Service, error) {
	var index bleve.Index
	loaded := false
	if _, err := os.Stat(indexPath); errors.Is(err, os.ErrNotExist) {
		log.Debug().Str("path", indexPath).Msg("creating index")
		err = os.MkdirAll(indexPath, os.FileMode(rwRwRw))
		if err != nil {
			return nil, fmt.Errorf("error creating index: %w", err)
		}
		contentMapping := bleve.NewDocumentMapping()
		contentMapping.Dynamic = false
		//// source data store - this is where original doc will be stored
		dataTextFieldMapping := bleve.NewTextFieldMapping()
		dataTextFieldMapping.Store = true
		dataTextFieldMapping.Index = true
		dataTextFieldMapping.IncludeInAll = false
		dataTextFieldMapping.IncludeTermVectors = true
		dataTextFieldMapping.IncludeInAll = false
		contentMapping.AddFieldMappingsAt("data", dataTextFieldMapping)
		englishTextFieldMapping := bleve.NewTextFieldMapping()
		englishTextFieldMapping.Analyzer = en.AnalyzerName
		contentMapping.AddFieldMappingsAt("match", englishTextFieldMapping)
		//// create
		indexMapping := bleve.NewIndexMapping()
		indexMapping.AddDocumentMapping("content", contentMapping)
		indexMapping.TypeField = "type"
		indexMapping.DefaultAnalyzer = "en"
		index, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			return nil, fmt.Errorf("error creating search index: %w", err)
		}
	} else {
		log.Debug().Str("path", indexPath).Msg("loading index")
		index, err = bleve.Open(indexPath)
		if err != nil {
			return nil, fmt.Errorf("error opening search index: %w", err)
		}
		loaded = true
	}
	return &search{index, loaded}, nil
}
