package psrozklad

import (
	"net/http"
	"testing"
)

type MockHttpClient struct {
	resp http.Response
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return &m.resp, nil
}

func TestNew(t *testing.T) {
	baseUri := "https://dekanat.zu.edu.ua/"
	want := Api{BaseUri: baseUri + "cgi-bin/timetable_export.cgi?&req_format=json&coding_mode=UTF8", HttpClient: http.DefaultClient}
	got := New(baseUri)
	if got != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
