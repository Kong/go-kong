package kong

import (
	"bytes"
	"context"
	"encoding/json"
)

// ListOpt aids in paginating through list endpoints
type ListOpt struct {
	// Size of the page
	Size int `url:"size,omitempty"`
	// Offset for the current page
	Offset string `url:"offset,omitempty"`

	// Tags to use for filtering the list.
	Tags []*string `url:"tags,omitempty"`
	// Tags are ORed by default, meaning entities
	// containing even a single tag in the list are listed.
	// If true, tags are ANDed, meaning only entities
	// matching each tag in the Tags array are listed.
	MatchAllTags bool
}

// qs is used to construct query string for list endpoints
type qs struct {
	Size   int    `url:"size,omitempty"`
	Offset string `url:"offset,omitempty"`
	Tags   string `url:"tags,omitempty"`
}

//Instantiate a ListOpt with the default page size an a deduplicted list of tags when present
func newOpt(tags []string) *ListOpt {
	opt := new(ListOpt)
	opt.Size = pageSize
	opt.Tags = StringSlice(deduplicate(tags)...)
	opt.MatchAllTags = true
	return opt
}

// list fetches a list of an entity in Kong.
// opt can be used to control pagination and tags
// allowNotFound allow to return an empty list for entities that are disabled or just doesn't exists on the used version
func (c *Client) listAll(ctx context.Context, endpoint string, opt *ListOpt, allowNotFound bool) ([]json.RawMessage, error) {
	var list, data []json.RawMessage
	var err error

	for opt != nil {
		data, opt, err = c.list(ctx, endpoint, opt)
		if allowNotFound && IsNotFoundErr(err) {
			return list, nil
		}
		if err != nil {
			return nil, err
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		list = append(list, data...)
	}
	return list, nil
}

// list fetches a list of an entity in Kong.
// opt can be used to control pagination.
func (c *Client) list(ctx context.Context,
	endpoint string, opt *ListOpt) ([]json.RawMessage, *ListOpt, error) {

	q := constructQueryString(opt)
	req, err := c.NewRequest("GET", endpoint, &q, nil)
	if err != nil {
		return nil, nil, err
	}
	var list struct {
		Data []json.RawMessage `json:"data"`
		Next *string           `json:"offset"`
	}

	_, err = c.Do(ctx, req, &list)
	if err != nil {
		return nil, nil, err
	}

	// convinient for end user to use this opt till it's nil
	var next *ListOpt
	if list.Next != nil {
		next = &ListOpt{
			Offset: *list.Next,
		}
		if opt != nil && next != nil {
			next.Size = opt.Size
			next.Tags = opt.Tags
			next.MatchAllTags = opt.MatchAllTags
		}
	}

	return list.Data, next, nil
}

func constructQueryString(opt *ListOpt) qs {
	var q qs
	if opt == nil {
		return q
	}
	q.Size = opt.Size
	q.Offset = opt.Offset
	var tagQS bytes.Buffer
	tagCount := len(opt.Tags)
	for i := 0; i < tagCount; i++ {
		tagQS.WriteString(*opt.Tags[i])
		if i+1 < tagCount {
			if opt.MatchAllTags {
				tagQS.WriteByte(',')
			} else {
				tagQS.WriteByte('/')
			}
		}
	}
	q.Tags = tagQS.String()

	return q
}
