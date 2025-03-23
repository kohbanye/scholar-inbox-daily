package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/kohbanye/scholar-inbox-daily/internal/domain"
)

type ScholarScraper struct{}

func NewScholarScraper() *ScholarScraper {
	return &ScholarScraper{}
}

func (s *ScholarScraper) FetchPapers() ([]domain.Paper, error) {
	email := os.Getenv("SCHOLAR_INBOX_EMAIL")
	password := os.Getenv("SCHOLAR_INBOX_PASSWORD")

	if email == "" || password == "" {
		return nil, fmt.Errorf("SCHOLAR_INBOX_EMAIL and SCHOLAR_INBOX_PASSWORD environment variables must be set")
	}

	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\nError: %v", r.Request.URL, r, err)
	})

	loginSuccess := false
	var papers []domain.Paper

	c.OnResponse(func(r *colly.Response) {
		if r.Request.URL.String() == "https://api.scholar-inbox.com/api/password_login" {
			var loginResponse map[string]interface{}
			if err := json.Unmarshal(r.Body, &loginResponse); err != nil {
				log.Printf("Error parsing login response: %v", err)
				return
			}

			if success, ok := loginResponse["success"].(bool); ok && success {
				loginSuccess = true
				log.Println("Login successful")

				c.Visit("https://api.scholar-inbox.com/api/")
			} else {
				log.Printf("Login failed: %v", loginResponse)
			}
		} else if r.Request.URL.String() == "https://api.scholar-inbox.com/api/" && loginSuccess {
			var apiResponse map[string]interface{}
			if err := json.Unmarshal(r.Body, &apiResponse); err != nil {
				log.Printf("Error parsing API response: %v", err)
				return
			}

			digestPapers, ok := apiResponse["digest_df"].([]interface{})
			if !ok {
				log.Printf("No digest_df field found in API response")
				return
			}

			for _, p := range digestPapers {
				paperMap, ok := p.(map[string]interface{})
				if !ok {
					continue
				}

				title, _ := paperMap["title"].(string)
				if title == "" {
					continue // Skip papers without a title
				}

				abstract, _ := paperMap["abstract"].(string)

				var authors string
				if authorsValue, ok := paperMap["authors"].(string); ok {
					authors = authorsValue
				} else if shortenedAuthors, ok := paperMap["shortened_authors"].(string); ok {
					authors = shortenedAuthors
				}

				url, _ := paperMap["url"].(string)

				var publishDate string
				if pubDate, ok := paperMap["publication_date"].(string); ok {
					publishDate = pubDate
				} else if displayVenue, ok := paperMap["display_venue"].(string); ok {
					publishDate = displayVenue
				}

				paper := domain.Paper{
					Title:    title,
					Authors:  authors,
					Date:     publishDate,
					URL:      url,
					Abstract: abstract,
				}

				papers = append(papers, paper)
			}
		}
	})

	loginData := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling login data: %v", err)
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	})

	err = c.PostRaw("https://api.scholar-inbox.com/api/password_login", jsonData)
	if err != nil {
		return nil, fmt.Errorf("error sending login request: %v", err)
	}

	c.Wait()

	if !loginSuccess {
		return nil, fmt.Errorf("login failed")
	}

	return papers, nil
}
