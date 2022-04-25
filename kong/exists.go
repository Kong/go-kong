package kong

import (
	"context"
	"net/http"
)

// exists check the existence  with a HEAD HTTP verb
func (c *Client) exists(ctx context.Context,
	endpoint string,
) (bool, error) {
	// Originally, this used HEAD. We promptly discovered that HEAD doesn't actually work for this
	// in Kong <=2.6: https://github.com/Kong/kong/issues/7554
	// Although future versions will support HEAD for existence checks, using GET for backwards
	// compatibility at the cost of efficiency
	req, err := c.NewRequest("GET", endpoint, nil, nil)
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
