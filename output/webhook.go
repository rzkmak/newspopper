package output

import (
	"bytes"
	"net/http"
)

type Webhook struct {
	url string
}

func NewWebhook(url string) Webhook {
	return Webhook{url: url}
}

func (w Webhook) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", w.url, bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return 0, err
	}
	return len(p), nil
}
