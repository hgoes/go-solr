package solr

import (
	"errors"
	"fmt"
)

type SolrResponse struct {
	Status int `json:"status"`
	QTime  int `json:"qtime"`
	Params struct {
		Query  string `json:"q"`
		Indent string `json:"indent"`
		Wt     string `json:"wt"`
	} `json:"params"`
	Response       Response `json:"response"`
	NextCursorMark string   `json:"nextCursorMark"`
	Adds           Adds     `json:"adds"`
}

type Response struct {
	NumFound uint32                   `json:"numFound"`
	Start    int                      `json:"start"`
	Docs     []map[string]interface{} `json:"docs"`
}

func GetDocIdFromDoc(m map[string]interface{}) string {
	if v, ok := m["id"]; ok {
		return v.(string)
	}
	return ""
}

func GetVersionFromDoc(m map[string]interface{}) int {
	if v, ok := m["_version_"]; ok {
		switch v := v.(type) {
		case float64:
			return int(v)
		case int:
			return v
		}
	}

	return 0
}

type Adds map[string]int

type UpdateResponse struct {
	Response struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
		RF     int `json:"rf"`
		MinRF  int `json:"min_rf"`
	} `json:"responseHeader"`

	// <adds> is a weird return value. It mixes ids with versions in a single slice, e.g.
	// ["id1",1233144,"id2",4122243], to get only the ids call AddedIDs afterwards
	Adds  []interface{} `json:"adds"`

	Error struct {
		Metadata []string `json:"metadata"`
		Msg      string   `json:"msg"`
		Code     int      `json:"code"`
	}
}

func (r UpdateResponse) AddedIDs() ([]string, error) {
	if len(r.Adds)%2 != 0 {
		return nil, errors.New("unexpected value for adds, len(adds) is not a multiple of 2")
	}
	ids := make([]string, len(r.Adds)/2)
	for i := range r.Adds {
		if i%2 == 0 {
			val, ok := r.Adds[i].(string)
			if !ok {
				return ids, fmt.Errorf("not a string: %v (position: %d)", val, i)
			}
			ids[i/2] = val
		}
	}
	return ids, nil
}

type DeleteRequest struct {
	Delete []string `json:"delete"`
}
