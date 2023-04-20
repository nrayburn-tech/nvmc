package util

import (
	"errors"
	"io"
	"net/http"
)

func Download(url string, destHandle io.Writer) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "nvmc-"+VERSION)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		if _, err = io.Copy(destHandle, response.Body); err != nil {
			return err
		}
		return nil
	case 300:
		fallthrough
	case 302:
		fallthrough
	case 307:
		redirectUrl := response.Header.Get("Location")
		if len(redirectUrl) > 0 {
			return Download(redirectUrl, destHandle)
		}
		return errors.New("300, 302, and 307 status codes must have a Location header")
	default:

	}

	return nil
}
