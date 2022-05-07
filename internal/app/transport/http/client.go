package http

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func NewHTTPClient() HTTPClient {
	return HTTPClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type HTTPClient struct {
	httpClient *http.Client
}

func (c HTTPClient) GetData(url string) ([]byte, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			err = errClose
		}
	}()

	return ioutil.ReadAll(resp.Body)
}

func (c HTTPClient) DownloadFile(filepath string, url string) error {
	data, err := c.GetData(url)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer func() {
		errClose := f.Close()
		if errClose != nil {
			err = errClose
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
