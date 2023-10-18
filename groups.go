package psrozklad

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Group struct {
	Departament string
	Name        string
	Id          int
}

// ID returns the ID of the group.
func (g Group) ID() int {
	return g.Id
}

// Type returns the type of the group, which is always "group".
func (g Group) Type() string {
	return "group"
}

// This function returns a list of groups from the API.
func (a *Api) GetGroups() ([]Group, error) {

	type group struct {
		Name string `json:"name"`
		Id   string `json:"ID"`
	}
	type departament struct {
		Name   string  `json:"name"`
		Groups []group `json:"objects"`
	}
	type psrozkladExport struct {
		Departments []departament `json:"departments"`
		Code        string        `json:"code"`
	}
	type export struct {
		PsrozkladExport psrozkladExport `json:"psrozklad_export"`
	}

	// Construct the URL for the API request.
	url := a.BaseUri
	url += "&req_type=obj_list"
	url += "&req_mode=group"
	url += "&show_ID=yes"

	// Create a new HTTP request.
	req, _ := http.NewRequest("GET", url, nil)

	// Execute the HTTP request.
	resp, err := a.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %v", err)
	}
	defer resp.Body.Close()

	// Decode the JSON response into an export struct.
	var exp export
	err = json.NewDecoder(resp.Body).Decode(&exp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode groups json: %v", err)
	}

	// Create a slice of Group structs to store the results.
	var groups []Group

	// Iterate over the departments in the export struct and add their groups to the results slice.
	for _, departament := range exp.PsrozkladExport.Departments {
		for _, group := range departament.Groups {
			// Convert the group ID to an integer.
			id, err := strconv.Atoi(group.Id)
			if err != nil {
				return nil, fmt.Errorf("failed to convert group id to string: %v", err)
			}

			// Add the group to the results slice.
			groups = append(groups, Group{
				Departament: departament.Name,
				Name:        group.Name,
				Id:          id,
			})
		}
	}

	// Return the results slice.
	return groups, nil
}
