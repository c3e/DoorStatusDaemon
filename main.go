package main

import (
	"goji.io"
	"goji.io/pat"
	"net/http"
	"github.com/spaceapi-community/go-spaceapi-spec/v13"
	"strconv"
	"os"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"log"
)

var url string

func main() {
	url = os.Getenv("API_URL") + "/api"

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/:location/:value"), setDoorAndStatus)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func setDoorAndStatus(w http.ResponseWriter, r *http.Request) {
	door := pat.Param(r, "location")
	status, err := strconv.ParseBool(pat.Param(r, "value"))

	if err != nil {
		http.Error(w, "couldn't parse " + pat.Param(r, "status") + " as bool", http.StatusBadRequest)
		return
	}

	err = setDoor(door, status, r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "couldn't send request", http.StatusInternalServerError)
		return
	}

	err = updateState(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "door set, however, couldn't set the correct state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}


func setDoor(location string, value bool, token string) error {
	foo := spaceapiStruct.DoorLocked{
		Value: value,
	}

	putData, err := json.Marshal(foo)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", url + "/sensors/door_locked/" + location, bytes.NewReader(putData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func updateState(token string) error {
	currentSpaceApi, err := loadSpaceApi()

	shouldOpen := false
	for _, doorLocked := range currentSpaceApi.Sensors.DoorLocked {
		if doorLocked.Value == false {
			shouldOpen = true
		}
	}

	if currentSpaceApi.State.Open == shouldOpen {
		return nil
	}

	newState := currentSpaceApi.State
	newState.Open = shouldOpen

	putData, err := json.Marshal(newState)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", url + "/api/state", bytes.NewReader(putData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func loadSpaceApi() (spaceapiStruct.SpaceAPI013, error) {
	resp, err := http.Get(url)
	if err != nil {
		return spaceapiStruct.SpaceAPI013{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return spaceapiStruct.SpaceAPI013{}, err
	}

	var currentSpaceApi spaceapiStruct.SpaceAPI013
	err = json.Unmarshal(body, &currentSpaceApi)
	if err != nil {
		return spaceapiStruct.SpaceAPI013{}, err
	}

	return currentSpaceApi, nil
}