package psrozklad

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Room struct {
	Block    string
	Name     string
	FullName string
	Id       int
}

// ID returns the ID of the room.
func (r Room) ID() int {
	return r.Id
}

// type_obj returns the type of the room, which is always "room".
func (r Room) type_obj() string {
	return "room"
}
// GetRooms gets list of all rooms.
func (a *Api) GetRooms() ([]Room, error) {
    // This function gets a list of rooms from the PS Rozklad API.

    type roomExport struct {
        Name string `json:"name"`
        Id   string `json:"ID"`
    }

    type blockExport struct {
        Name  string       `json:"name"`
        Rooms []roomExport `json:"objects"`
    }

    type psrozkladExport struct {
        Blocks []blockExport `json:"blocks"`
        Code   string        `json:"code"`
    }

    type export struct {
        PsrozkladExport psrozkladExport `json:"psrozklad_export"`
    }

    // Create a URL to the PS Rozklad API.
    url := a.BaseUri
    url += "&req_type=obj_list"
    url += "&req_mode=room"
    url += "&show_ID=yes"

    // Create a new HTTP request.
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create http request: %v", err)
    }

    // Do the HTTP request.
    resp, err := a.HttpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to do http request: %v", err)
    }
    defer resp.Body.Close()

    // Decode the JSON response into an export struct.
    var exp export
    err = json.NewDecoder(resp.Body).Decode(&exp)
    if err != nil {
        return nil, fmt.Errorf("failed to decode room list json: %v", err)
    }

    // Create a list of Room structs.
    var rooms []Room
    for _, block := range exp.PsrozkladExport.Blocks {
        for _, room := range block.Rooms {
            // Convert the ID string to an integer.
            id, err := strconv.Atoi(room.Id)
            if err != nil {
                return nil, fmt.Errorf("failed to convert id: %v, to string: %v", room.Id, err)
            }

            // Split the room name into two parts: block and name.
            n := strings.Split(room.Name, "/")

            // Add the room to the list of rooms.
            rooms = append(rooms, Room{
                Block:    block.Name,
                Name:     n[0],
                FullName: room.Name,
                Id:       id,
            })
        }
    }

    // Return the list of rooms.
    return rooms, nil
}
