package psrozklad

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Teacher struct {
	ShortName   string
	P           string
	I           string
	B           string
	Departament string
	Id          int
}

// This function returns the teacher's unique identifier.
func (t Teacher) ID() int {
	return t.Id
}

// This function returns the teacher's type, which is always "teacher".
func (t Teacher) Type() string {
	return "teacher"
}

// GetTeachers returns a list of teachers from the API.
func (a *Api) GetTeachers() ([]Teacher, error) {
	// Define a struct to represent a teacher export from the API.
	type teacherExport struct {
		Name string `json:"name"`
		P    string `json:"P"`
		I    string `json:"I"`
		B    string `json:"B"`
		Id   string `json:"ID"`
	}

	// Define a struct to represent a department export from the API.
	type departament struct {
		Name     string          `json:"name"`
		Teachers []teacherExport `json:"objects"`
	}

	// Define a struct to represent a psrozklad export from the API.
	type psrozkladExport struct {
		Departments []departament `json:"departments"`
		Code        string        `json:"code"`
	}

	// Define a struct to represent the overall export from the API.
	type export struct {
		PsrozkladExport psrozkladExport `json:"psrozklad_export"`
	}

	// Construct the API request URL.
	url := a.BaseUri
	url += "&req_type=obj_list"
	url += "&req_mode=teacher"
	url += "&show_ID=yes"

	// Create a new HTTP request.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %v", err)
	}

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
		return nil, fmt.Errorf("failed to decode teachers list json: %v", err)
	}

	// Create a slice of Teacher objects to store the results.
	var teachers []Teacher

	// Iterate over the departments and teachers in the export struct.
	for _, dep := range exp.PsrozkladExport.Departments {
		for _, teacher := range dep.Teachers {
			// Convert the teacher ID to an integer.
			id, err := strconv.Atoi(teacher.Id)
			if err != nil {
				return nil, fmt.Errorf("failed to convert id: %v, to string: %v", teacher.Id, err)
			}

			// Create a new Teacher object and append it to the slice.
			teachers = append(teachers, Teacher{
				ShortName:   teacher.Name,
				P:           teacher.P,
				I:           teacher.I,
				B:           teacher.B,
				Departament: dep.Name,
				Id:          id,
			})
		}
	}

	// Return the slice of Teacher objects.
	return teachers, nil
}
