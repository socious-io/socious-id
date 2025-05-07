package tests_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"socious-id/src/apps/shortener"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func shortenerGroup() {
	var shortLinks []*shortener.ShortenerURL
	It("should create shortener", func() {
		for _, url := range shorteningURLs {
			s, err := shortener.New(url)
			if err != nil {
				log.Fatalf("Shortener error : %v", err)
			}
			shortLinks = append(shortLinks, s)
		}
	})

	It("should fetch shortener", func() {
		for _, short := range shortLinks {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/fetch", short.ShortID), nil)
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(200))
			Expect(body["long_url"].(string)).To(Equal(short.LongURL))
		}
	})

	It("should  fetch and redirect shortener", func() {
		for _, short := range shortLinks {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", short.ShortID), nil)
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusSeeOther))
		}
	})
}
