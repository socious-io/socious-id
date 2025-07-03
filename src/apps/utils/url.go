package utils

import "net/url"

func CreateUrl(baseUrl string, query map[string]string) (string, error) {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	// Add query parameters
	params := url.Values{}
	for key, value := range query {
		params.Add(key, value)

	}
	base.RawQuery = params.Encode()

	return base.String(), nil
}
