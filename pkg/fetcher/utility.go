package fetcher

import (
	"context"
	"net/http"
)

func Request(ctx context.Context, targetURL string) (*http.Response, error) {
	// set 10 seconds timeout
	// use new context due to multiple calls
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
