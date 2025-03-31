package scraper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/kohbanye/scholar-inbox-daily/internal/domain"
	"github.com/kohbanye/scholar-inbox-daily/internal/logger"
)

type ScholarScraper struct {
	logger *logger.Logger
}

func NewScholarScraper() *ScholarScraper {
	return &ScholarScraper{
		logger: logger.New(),
	}
}

func (s *ScholarScraper) FetchPapers() ([]domain.Paper, error) {
	s.logger.Info("Starting to fetch papers from Scholar Inbox")
	email := os.Getenv("SCHOLAR_INBOX_EMAIL")
	password := os.Getenv("SCHOLAR_INBOX_PASSWORD")

	if email == "" || password == "" {
		err := fmt.Errorf("SCHOLAR_INBOX_EMAIL and SCHOLAR_INBOX_PASSWORD environment variables must be set")
		s.logger.Error("Environment variables not set", err)
		return nil, err
	}

	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		s.logger.Error(fmt.Sprintf("Request URL: %s failed", r.Request.URL), err)
	})

	loginSuccess := false
	var papers []domain.Paper

	c.OnResponse(func(r *colly.Response) {
		if r.Request.URL.String() == "https://api.scholar-inbox.com/api/password_login" {
			var loginResponse map[string]interface{}
			if err := json.Unmarshal(r.Body, &loginResponse); err != nil {
				s.logger.Error("Error parsing login response", err)
				return
			}

			if success, ok := loginResponse["success"].(bool); ok && success {
				loginSuccess = true
				s.logger.Info("Login successful")

				c.Visit("https://api.scholar-inbox.com/api/")
			} else {
				s.logger.Error("Login failed", fmt.Errorf("login response: %v", loginResponse))
			}
		} else if r.Request.URL.String() == "https://api.scholar-inbox.com/api/" && loginSuccess {
			var apiResponse map[string]interface{}
			if err := json.Unmarshal(r.Body, &apiResponse); err != nil {
				s.logger.Error("Error parsing API response", err)
				return
			}

			digestPapers, ok := apiResponse["digest_df"].([]interface{})
			if !ok {
				s.logger.Error("No digest_df field found in API response", fmt.Errorf("api response: %v", apiResponse))
				return
			}

			s.logger.Info(fmt.Sprintf("Found %d papers in API response", len(digestPapers)))

			for _, p := range digestPapers {
				paperMap, ok := p.(map[string]interface{})
				if !ok {
					continue
				}

				title, _ := paperMap["title"].(string)
				if title == "" {
					s.logger.Debug("Skipping paper without title")
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
		s.logger.Error("Error marshaling login data", err)
		return nil, fmt.Errorf("error marshaling login data: %v", err)
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	})

	s.logger.Info("Sending login request")
	err = c.PostRaw("https://api.scholar-inbox.com/api/password_login", jsonData)
	if err != nil {
		s.logger.Error("Error sending login request", err)
		return nil, fmt.Errorf("error sending login request: %v", err)
	}

	c.Wait()

	if !loginSuccess {
		err := fmt.Errorf("login failed")
		s.logger.Error("Login failed", err)
		return nil, err
	}

	s.logger.Info(fmt.Sprintf("Successfully fetched %d papers", len(papers)))
	return papers, nil
}
