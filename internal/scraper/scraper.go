package scraper

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gocolly/colly/v2"
)

func ScrapeShipCode(imoNumber int, userAgent string) (string, error) {
	url := fmt.Sprintf("https://www.ship-db.de/SuNamIMO.php?Sprache=E&suche=%d&Exakt=Exakt&senden=Name+%%2F+IMO", imoNumber)

	// Almacenamiento para los datos de las etiquetas
	var rawHTML string

	// Leer el certificado
	caCert, err := os.ReadFile("./shipdb.pem")
	if err != nil {
		log.Fatalf("Error al leer el certificado CA: %v", err)
	}

	// Crear un pool de certificados raíz
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configurar un transporte HTTP con los certificados raíz personalizados
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: caCertPool},
	}

	// Inicializar Colly
	c := colly.NewCollector(colly.UserAgent(userAgent))
	c.SetClient(&http.Client{Transport: tr})

	// Manejar errores y capturar el cuerpo de la respuesta
	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 500 && len(r.Body) > 0 {
			// Procesar el HTML del cuerpo para extraer los datos
			rawHTML = string(r.Body)
		}
	})

	// Visitar la URL
	err = c.Visit(url)
	if err != nil {
		return rawHTML, fmt.Errorf("error StatusCode: %w", err)
	}

	return "", nil
}

func ScrapeData(url string, userAgent string) (string, error) {
	var result string

	caCert, err := os.ReadFile("./shipdb.pem")
	if err != nil {
		log.Fatalf("Error al leer el certificado CA: %v", err)
	}

	// Crear un pool de certificados raíz
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configurar un transporte HTTP con los certificados raíz personalizados
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: caCertPool},
	}

	c := colly.NewCollector(colly.UserAgent(userAgent))
	c.SetClient(&http.Client{Transport: tr})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Ajustar el selector al contenido que se necesita
		result = e.Text
	})

	c.OnError(func(r *colly.Response, err error) {
		// Manejar errores de scraping
		log.Printf("Errror scraping %s: %v", url, err)
	})

	err = c.Visit(url)
	if err != nil {
		return "", err
	}
	return result, nil
}
