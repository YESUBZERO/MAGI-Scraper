package parsership

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ParserShipData(data string) (string, error) {
	// Normalizar los datos
	rawData := strings.ReplaceAll(data, "\t", " ")
	rawData = strings.ReplaceAll(rawData, "\u00A0", " ") // Reemplazar espacios Unicode
	rawData = regexp.MustCompile(`\s+`).ReplaceAllString(rawData, " ")
	rawData = strings.TrimSpace(rawData)

	// Diccionario para almacenar los datos
	parsedData := make(map[string]string)

	// Lista ordenada de claves
	orderedKeys := []string{
		"Type",
		"Built",
		"IMO-No.",
		"Shipyard",
		"Hull-No.",
		"Keel Laying",
		"Launch",
		"Delivery",
		"GT",
		"NT",
		"Carrying capacity (tdw)",
		"Length overall (m)",
		"Breadth (m)",
		"Depth (m)",
		"Propulsion",
		"Power",
		"Screws",
		"Speed",
	}

	// Patrones de extracción ajustados
	patterns := map[string]string{
		"Type":                    `Type:\s*([^\n]+?)\s+Built:`,
		"Built":                   `Built:\s*(\d{4})`,
		"IMO-No.":                 `IMO-No.:\s*(\d+)`,
		"Shipyard":                `Shipyard:\s*([^\n]+?)\s+Hull-No.:`,
		"Hull-No.":                `Hull-No.:\s*([^\n]+?)\s+Keel Laying:`,
		"Keel Laying":             `Keel Laying:\s*([\d.]+)`,
		"Launch":                  `Launch:\s*([\d.]+)`,
		"Delivery":                `Delivery:\s*([\d.]+)`,
		"GT":                      `GT\s*([\d.,]+)`,
		"NT":                      `NT\s*([\d.,]+)`,
		"Carrying capacity (tdw)": `Carrying capacity \(tdw\):\s*([\d.,]+)`,
		"Length overall (m)":      `Length overall \(m\):\s*([\d.,]+)`,
		"Breadth (m)":             `Breadth \(m\):\s*([\d.,]+)`,
		"Depth (m)":               `Depth \(m\):\s*([\d.,]+)`,
		"Propulsion":              `Propulsion:\s*([^\n]+?)\s+Power:`,
		"Power":                   `(?i)Power\s*:\s*([\d.,]+)\s*(KW|HP)`,
		"Screws":                  `Screws:\s*(\d+)`,
		"Speed":                   `Speed:\s*([\d.,]+)\s*Knot`,
	}

	// Aplicar patrones
	for key, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(rawData)
		if len(match) > 1 {
			if key == "Power" {
				// Manejar Power con unidad (KW o HP)
				value := match[1]
				unit := strings.ToUpper(strings.TrimSpace(match[2]))
				if unit == "HP" {
					// Convertir HP a KW
					parsedValue, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
					if err == nil {
						parsedData[key] = fmt.Sprintf("%.2f KW", parsedValue*0.7355)
					}
				} else {
					parsedData[key] = fmt.Sprintf("%s %s", value, unit)
				}
			} else {
				parsedData[key] = strings.TrimSpace(match[1])
			}
		}
	}

	// Construir JSON ordenado manualmente
	var jsonBuilder strings.Builder
	jsonBuilder.WriteString("{\n")
	for i, key := range orderedKeys {
		if value, exists := parsedData[key]; exists {
			jsonBuilder.WriteString(fmt.Sprintf("  \"%s\": \"%s\"", key, value))
			if i < len(orderedKeys)-1 {
				jsonBuilder.WriteString(",\n")
			} else {
				jsonBuilder.WriteString("\n")
			}
		}
	}
	jsonBuilder.WriteString("}")

	return jsonBuilder.String(), nil
}
