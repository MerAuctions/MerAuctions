package api

// package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/MerAuctions/MerAuctions/models"
)

func GetTagsForImage(name string) models.TagList {
	var tagList models.TagList
	var w http.ResponseWriter
	var IP = "34.84.23.245"
	imageURL := "http://" + IP + "/images/" + name
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

func GetDescriptionForImage(name string) models.Description {
	var description models.Description
	var w http.ResponseWriter
	var IP = "34.84.23.245"
	imageURL := "http://" + IP + "/images/" + name
	resp, err := http.Get("https://cat-that-pic-bak.herokuapp.com/api/v1/getcaption?fileName=" + imageURL)
	if err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return description
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return description
	}

	description = models.Description(body)

	return description
}
