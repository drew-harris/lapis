package mojang

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UUIDSuccessResponse struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func GetUUID(username string) (string, error) {
	fmt.Println("Getting UUID for " + username + " from mojang servers...")

	req, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + username)
	if err != nil {
		return "", err
	}

	if req.StatusCode == 200 {
		body := UUIDSuccessResponse{}
		json.NewDecoder(req.Body).Decode(&body)
		req.Body.Close()
		fmt.Println("Got UUID for " + username + " from mojang servers")
		return body.ID, nil
	}
	if req.StatusCode == 400 {
		return "", errors.New("Invalid Minecraft Username")
	}

	return "", errors.New("Couldn't get uuid from mojang servers")
}
