package todoist

import "net/http"
import "fmt"
import "io/ioutil"

const DefaultRestUrl string = "https://beta.todoist.com/API/v8"
const DefaultSyncUrl string = "https://todoist.com/api/v7"

type Client struct {
	HttpClient *http.Client
	ApiKey     string
	Host       string
	Base       string
}

func NewClient(apiKey string) *Client {
	return &Client{
		HttpClient: http.DefaultClient,
		ApiKey:     apiKey,
	}
}

func (c *Client) newRequest(path string) (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", DefaultRestUrl, path), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent {
		return body, err
	} else {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}
}
