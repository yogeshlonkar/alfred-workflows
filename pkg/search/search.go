package search

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/rs/zerolog/log"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
)

type search struct {
	bleve.Index
	loaded bool
}

func (s *search) Load(suggestion ...alfred.Item) error {
	log.Debug().Int("count", len(suggestion)).Msg("adding suggestion to index")
	for _, suggestion := range suggestion {
		data, err := json.Marshal(suggestion)
		if err != nil {
			return fmt.Errorf("error marshalling suggestion: %w", err)
		}
		// log.Debug().Str("data", string(data)).Msg("adding")
		err = s.Index.Index(suggestion.UID, SuggestionIndexMeta{suggestion.UID, string(data), suggestion.Match})
		if err != nil {
			return fmt.Errorf("error indexing suggestion: %w", err)
		}
	}
	return nil
}

func (s *search) Match(keyword string) ([]alfred.Item, error) {
	log.Debug().Str("keyword", keyword).Msg("searching in index")
	query := bleve.NewDisjunctionQuery()
	keywords := strings.Split(strings.ToLower(keyword), " ")
	for index, term := range keywords {
		log.Debug().Str("term", term).Msg("added FuzzyQuery")
		pQuery := bleve.NewPrefixQuery(term)
		pQuery.SetField("match")
		pQuery.SetBoost(float64(len(keywords) - index - 1))
		query.AddQuery(pQuery)
		// tQuery := bleve.NewTermQuery(term)
		// tQuery.SetField("match")
		// tQuery.SetBoost(float64(len(keywords) - index))
		// query.AddQuery(tQuery)
	}
	req := bleve.NewSearchRequest(query)
	req.Fields = []string{"data"}
	req.SetSortFunc(func(s sort.Interface) {
	})
	// req.Explain = true
	searchResults, err := s.Search(req)
	if err != nil {
		return nil, fmt.Errorf("error searching suggestions: %w", err)
	}
	log.Debug().Int("hits", searchResults.Hits.Len()).Msg("found")
	suggestions := make([]alfred.Item, 0)
	for _, hit := range searchResults.Hits {
		var suggestion alfred.Item
		err = json.Unmarshal([]byte(hit.Fields["data"].(string)), &suggestion)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling suggestion: %w", err)
		}
		// log.Debug().Str("title", suggestion.Title).Interface("hit.Expl", hit.Expl).Send()
		suggestions = append(suggestions, suggestion)
	}
	return suggestions, nil
}

func (s *search) Close() error {
	err := s.Index.Close()
	if err != nil {
		return fmt.Errorf("error closing index: %w", err)
	}
	return nil
}

func (s *search) Loaded() bool {
	return s.loaded
}
