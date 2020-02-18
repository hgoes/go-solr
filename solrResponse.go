package solr

import (
	"encoding/json"
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
	Debug          *Debug   `json:"debug"`
	NextCursorMark string   `json:"nextCursorMark"`
	Adds           Adds     `json:"adds"`
}

type Debug struct {
	Timing Timing `json:"timing"`
}

type Timing struct {
	Time    uint32      `json:"time"`
	Prepare StageTiming `json:"prepare"`
	Process StageTiming `json:"process"`
}

type StageTiming struct {
	Time uint32 `json:"time"`
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
	Adds []interface{} `json:"adds"`

	Error struct {
		Metadata []string `json:"metadata"`
		Msg      string   `json:"msg"`
		Code     int      `json:"code"`
	}
}

func (r UpdateResponse) AddedIDs() (ids map[string]int64, err error) {
	if len(r.Adds)%2 != 0 {
		return nil, errors.New("unexpected value for adds, len(adds) is not a multiple of 2")
	}
	ids = make(map[string]int64, len(r.Adds)/2)
	for i := range r.Adds {
		if i%2 == 0 {
			id, ok := r.Adds[i].(string)
			if !ok {
				return ids, fmt.Errorf("not a string: %v (position: %d)", id, i)
			}
			rev, ok := r.Adds[i+1].(json.Number)
			if !ok {
				return ids, fmt.Errorf("not a revision: %v (position: %d)", rev, i+1)
			}
			ids[id], err = rev.Int64()
			if err != nil {
				return ids, fmt.Errorf("not a revision: %v (position: %d)", rev, i+1)
			}
		}
	}
	return ids, nil
}

type DeleteRequest struct {
	Delete []string `json:"delete"`
}
