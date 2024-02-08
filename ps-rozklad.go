// Package psrozklad providae api wrapper for ps-rozklad api in golang
package psrozklad

import (
	"fmt"
	"net/http"
	"strings"
)

type Api struct {
	BaseUri    string
	HttpClient Client
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Rooms    map[string]Room
	Teachers map[string]Teacher
	Groups   map[string]Group
)

// New creates a new Api struct.
func New(BaseUri string) Api {
	// Add some query parameters to the base URI.
	BaseUri += "cgi-bin/timetable_export.cgi?&req_format=json&coding_mode=UTF8"

	// Create a new Api struct and return it.
	return Api{BaseUri: BaseUri, HttpClient: http.DefaultClient}
}

// Initialize rooms,groups,teschers
func (a *Api) Init() error {
	// Initialize the groups.
	err := a.InitGroups()
	if err != nil {
		// Return an error if the groups failed to initialize.
		return fmt.Errorf("failed to init: %v", err)
	}

	// Initialize the rooms.
	err = a.InitRooms()
	if err != nil {
		// Return an error if the rooms failed to initialize.
		return fmt.Errorf("failed to init: %v", err)
	}

	// Initialize the teachers.
	err = a.InitTeachers()
	if err != nil {
		// Return an error if the teachers failed to initialize.
		return fmt.Errorf("failed to init: %v", err)
	}

	// Return nil if all of the initialization functions succeeded.
	return nil
}

// InitRooms initializes the Rooms map.
func (a *Api) InitRooms() error {
	// Create a new empty map.
	Rooms = make(map[string]Room)

	// Get a list of all of the rooms from the API.
	rooms, err := a.GetRooms()
	if err != nil {
		// Return an error if the GetRooms() function failed.
		return fmt.Errorf("failed to get rooms: %v", err)
	}

	// Iterate over the list of rooms and add each room to the Rooms map.
	for _, room := range rooms {
		// The key of each entry in the map is the room's name and block, concatenated with a slash.
		key := room.Name + "/" + room.Block
		Rooms[strings.ToLower(key)] = room
	}

	// Return nil, indicating that the function was successful.
	return nil
}

// InitGroups initializes the Groups map.
func (a *Api) InitGroups() error {
	// Create a new empty map.
	Groups = make(map[string]Group)

	// Get a list of all of the groups from the API.
	groups, err := a.GetGroups()
	if err != nil {
		// Return an error if the GetGroups() function failed.
		return fmt.Errorf("failed to get groups: %v", err)
	}

	// Iterate over the list of groups and add each group to the Groups map.
	for _, group := range groups {
		// The key of each entry in the map is the group's name.
		Groups[strings.ToLower(group.Name)] = group
	}

	// Return nil, indicating that the function was successful.
	return nil
}

// InitTeachers initializes the Teachers map.
func (a *Api) InitTeachers() error {
	// Create a new empty map.
	Teachers = make(map[string]Teacher)

	// Get a list of all of the teachers from the API.
	teachers, err := a.GetTeachers()
	if err != nil {
		// Return an error if the GetTeachers() function failed.
		return fmt.Errorf("failed to get teachers: %v", err)
	}

	// Iterate over the list of teachers and add each teacher to the Teachers map.
	for _, teacher := range teachers {
		// The key of each entry in the map is the teacher's short name.
		Teachers[strings.ToLower(teacher.ShortName)] = teacher
	}

	// Return nil, indicating that the function was successful.
	return nil
}
