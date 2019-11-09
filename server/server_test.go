package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/MerAuctions/MerAuctions/server"
)

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

var _ = Describe("Server", func() {
	var (
		router   *gin.Engine
		response *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		router = CreateRouter()
		// Since we modify lists in memory, we need to restore them to a clean state before every test
	})

	Describe("Version 1 API at /api/v0.1", func() {
		Describe("The / endpoint", func() {
			BeforeEach(func() {
				response = performRequest(router, "GET", "/auctions/1", nil)
			})

			It("Returns with Status 200", func() {
				Expect(response.Code).To(Equal(200))
			})

			// It("Returns the String 'Hello World'", func() {
			// 	Expect(response.Body.String()).To(Equal("Hello World"))
			// })

		})
	})
})
