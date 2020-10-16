package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func GitHubApprove(caller, pr, username, password string) error {
	url, err := prIntoReviewAddress(pr, username, password)
	if err != nil {
		return err
	}

	body := fmt.Sprintf(`{"body":"%s 🤖 👍","event":"APPROVE"}`, caller)
	client := http.Client{}
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s\n%s", resp.Status, respBody)
	}
	return nil
}

var rPR = regexp.MustCompile(`https://github.com/(.*)/(.*)/pull/(.*)`)

func prIntoReviewAddress(pr, username, password string) (*url.URL, error) {
	matches := rPR.FindStringSubmatch(pr)
	if len(matches) == 0 {
		return nil, fmt.Errorf("not a valid pr: %s", pr)
	}

	u, _ := url.Parse(fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%s/reviews", matches[1], matches[2], matches[3]))
	u.User = url.UserPassword(username, password)
	return u, nil
}
