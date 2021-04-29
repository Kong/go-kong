package kong

import (
	"context"
)

// AbstractTagService handles Tags in Kong.
type AbstractTagService interface {
	//Exists checks if the tags exists
	Exists(ctx context.Context) (*bool, error)
}

// TagService handles Tags in Kong.
type TagService service

// Exists checks exitence of the Tags in Kong.
func (s *TagService) Exists(ctx context.Context) (*bool, error) {
	return s.client.exists(ctx, "/tags")
}
