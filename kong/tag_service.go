package kong

import (
	"context"
	"net/http"
)

// AbstractTagService handles Tags in Kong.
type AbstractTagService interface {
	// Check if the tags exists
	Exists(ctx context.Context) (*bool, error)
}

// TagService handles Tags in Kong.
type TagService service

// Exists check exitence of the Tags in Kong.
func (s *TagService) Exists(ctx context.Context) (*bool, error) {
	req, err := s.client.NewRequest("HEAD", "/tags", nil, nil)
	if err != nil {
		return nil, err
	}
	var status = false
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		if IsNotFoundErr(err) {
			return &status, nil
		}
		return nil, err
	}
	status = resp.StatusCode == http.StatusOK
	return &status, nil
}
