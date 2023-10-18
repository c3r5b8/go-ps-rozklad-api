package psrozklad_test

import (
	"net/http"
	"testing"

	psrozklad "github.com/c3r5b8/go-ps-rozklad-api"
)

type MockHttpClient struct {
	resp http.Response
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return &m.resp, nil
}

func TestNew(t *testing.T) {
	baseUri := "https://dekanat.zu.edu.ua/"
	want := psrozklad.Api{BaseUri: baseUri + "cgi-bin/timetable_export.cgi?&req_format=json&coding_mode=UTF8", HttpClient: http.DefaultClient}
	got := psrozklad.New(baseUri)
	if got != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
