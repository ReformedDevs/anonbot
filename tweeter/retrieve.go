package tweeter

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
)

// retrieveMedia fetches a remote resource and converts its contents to base64.
func (t *Tweeter) retrieveMedia(mediaURL string) (string, error) {
	t.log.Infof("retrieving %s...", mediaURL)
	req, err := http.NewRequest(http.MethodGet, mediaURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "anonbot - https://git.io/vdzBx")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	var (
		b = &bytes.Buffer{}
		w = base64.NewEncoder(base64.StdEncoding, b)
	)
	if _, err := io.Copy(w, res.Body); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return b.String(), nil
}
