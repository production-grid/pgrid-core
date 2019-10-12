package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/production-grid/pgrid-core/pkg/applications"
)

//PostJSON is used to test unauthenticated post requests as JSON
func PostJSON(t *testing.T, pathInfo string, body interface{}, responsePtr interface{}) {

	httpClient := &http.Client{}

	content, err := json.Marshal(body)
	HandlePossibleError(t, err)

	path := TestServer.URL + pathInfo

	req, err := http.NewRequest("POST", path, bytes.NewBuffer(content))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	HandlePossibleError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Error(resp.Status)
	}

	ConsumeResponse(responsePtr, resp)

}

//GetJSONWithSessionKey retrieves json with a given session token
func GetJSONWithSessionKey(t *testing.T, pathInfo string, sessionKey string, responsePtr interface{}) {

	httpClient := &http.Client{}

	path := TestServer.URL + pathInfo

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if sessionKey != "" {
		req.AddCookie(&http.Cookie{Name: applications.SessionCookieName, Path: "/", Value: sessionKey})
	}
	resp, err := httpClient.Do(req)
	HandlePossibleError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Error(resp.Status)
	}

	ConsumeResponse(responsePtr, resp)

}

/*
ConsumeResponse parses a JSON response and uses it to populate the given entity
interface.
*/
func ConsumeResponse(responseEntity interface{}, resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	fmt.Println("Response JSON:", string(b))

	err = json.Unmarshal(b, responseEntity)

	if err != nil {
		return err
	}

	return nil

}

/*
HandlePossibleError encapsulates a if statement that might
become monotonous for test developers.
*/
func HandlePossibleError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
