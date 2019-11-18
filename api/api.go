package api

// package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/MerAuctions/MerAuctions/models"
)

func GetTagsForImage(name string) models.TagList {
	var tagList models.TagList
	var w http.ResponseWriter
	imageURL := "http://" + string(os.Getenv("DOMAIN")) + "/images/" + name
	resp, err := http.Get("https://cat-that-pic-bak.herokuapp.com/api/v1/getTagsfromImage?fileName=" + imageURL)
	if err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return tagList
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&tagList); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return tagList
	}

	return tagList
}
