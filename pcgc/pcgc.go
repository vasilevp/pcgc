package pcgc

import (
	"context"
	"fmt"
	"net/http"
)

// DoGetRequest creates a new GET HTTP request, executes it, and returns the response
func DoGetRequest(url string) (*http.Response, error) {
	ctx, onCancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer onCancel()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}

	return http.DefaultClient.Do(req.WithContext(ctx))
}
