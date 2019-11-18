package api

// package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"

	"github.com/MerAuctions/MerAuctions/models"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetTagsForImage(name string) models.TagList {
	var tagList models.TagList
	var w http.ResponseWriter
	imageURL := "http://" + string(os.Getenv("DOMAIN")) + "/images/" + name

	//imageURL := "https://cms.hostelworld.com/hwblog/wp-content/uploads/sites/2/2017/08/lovelyforliving.jpg"
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
