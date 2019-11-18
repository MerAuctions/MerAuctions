package api

// package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/MerAuctions/MerAuctions/models"
)

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

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
	imageURL := getLocalIP() + "/images/" + name
	response := performRequest(router, "GET", imageURL, nil)
	json.Unmarshal(response.Body.Bytes(), &tagList)
	return tagList
}
