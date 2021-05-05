package kong

import (
	"context"
	"net/http"
)

// exists check the existence  with a HEAD HTTP verb
func (c *Client) exists(ctx context.Context,
	endpoint string) (bool, error) {
	req, err := c.NewRequest("HEAD", endpoint, nil, nil)
	if err != nil {
		return false, err
	}
	resp, err := c.Do(ctx, req, nil)
	if err != nil {
		if IsNotFoundErr(err) {
			return false, nil
		}
		return false, err
	}
	return resp.StatusCode == http.StatusOK, nil
}
