package parsership

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParserShipCode(rawHTML string) (string, error) {
	// Procesar el HTML del cuerpo para extraer los datos
	var shipURL string
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(rawHTML)))
	if err != nil {
		return "", fmt.Errorf("error al crear el documento goquery: %w", err)
	}

	// Encontrar la URL del barco
	doc.Find("td a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if exists {
			shipURL = href
			return false
		}
		return true
	})

	if shipURL == "" {
		return "", fmt.Errorf("no se encontró la URL del barco")
	}

	// Reemplazar espacios en la URL
	shipURL = strings.ReplaceAll(shipURL, " ", "%20")

	parsedURL, err := url.Parse(shipURL)
	if err != nil {
		return "", fmt.Errorf("error al analizar la URL extraída: %w", err)
	}

	shipURL = parsedURL.String()

	return shipURL, nil

}
