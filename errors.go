package solr

import (
	"fmt"
	"net/http"
)

// HttpError - An Error due to a http response >= 400
type HttpError struct {
	Message string
	Status int
}

func (err HttpError) Error() string {
	return fmt.Sprintf("received error response from solr status: %d message: %s", err.Status, err.Message)
}

func isHttpNotFound(err error) bool {
	if httpErr, ok := err.(HttpError); ok && httpErr.Status == http.StatusNotFound {
		return true
	}
	return false
}

type SolrError struct {
	errorMessage string
}

func (err SolrError) Error() string {
	return err.errorMessage
}

func NewSolrError(status int, message string) error {
	return SolrError{errorMessage: fmt.Sprintf("received error response from solr status: %d message: %s", status, message)}
}

func NewSolrRFError(rf, minRF int) error {
	return SolrMinRFError{SolrError{errorMessage: fmt.Sprintf("received error response from solr: rf (%d) is < min_rf (%d)", rf, minRF)}, rf}
}

type SolrMinRFError struct {
	SolrError
	MinRF int
}

type SolrInternalError struct {
	SolrError
}

func NewSolrInternalError(status int, message string) error {
	return SolrInternalError{SolrError{errorMessage: fmt.Sprintf("received error response from solr status: %d message: %s", status, message)}}
}

type SolrLeaderError struct {
	SolrError
}

func NewSolrLeaderError(docID string) error {
	return SolrLeaderError{SolrError{errorMessage: fmt.Sprintf("Cannot find leader for doc %s", docID)}}
}

type SolrBatchError struct {
	error
}

func NewSolrBatchError(err error) error {
	return SolrBatchError{error: err}
}

type SolrParseError struct {
	SolrError
}

func NewSolrParseError(status int, message string) error {
	return SolrInternalError{SolrError{errorMessage: fmt.Sprintf("received error response from solr status: %d message: %s", status, message)}}
}

type SolrMapParseError struct {
	bucket string
	m      map[string]interface{}
	userId int
}

func (err SolrMapParseError) Error() string {
	return fmt.Sprintf("SolrMapParseErr: map does not contain email_register, bucket: %s, userId: %d map: %v", err.bucket, err.userId, err.m)

}
func NewSolrMapParseError(bucket string, userId int, m map[string]interface{}) error {
	return SolrMapParseError{bucket, m, userId}
}

