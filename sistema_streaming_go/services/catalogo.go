package services

import (
	"errors"
	"strings"

	"sistema-streaming/models"
)

// BuscarPorTitulo sigue un enfoque funcional: recibe datos y devuelve
// un nuevo resultado sin modificar el catálogo original.
func BuscarPorTitulo(catalogo []models.Reproducible, texto string) []models.Reproducible {
	texto = strings.ToLower(strings.TrimSpace(texto))
	resultado := make([]models.Reproducible, 0)

	for _, contenido := range catalogo {
		if strings.Contains(strings.ToLower(contenido.GetTitulo()), texto) {
			resultado = append(resultado, contenido)
		}
	}
	return resultado
}

// FiltrarPorGenero también devuelve un nuevo slice y no altera el original.
func FiltrarPorGenero(catalogo []models.Reproducible, genero string) []models.Reproducible {
	genero = strings.ToLower(strings.TrimSpace(genero))
	resultado := make([]models.Reproducible, 0)

	for _, contenido := range catalogo {
		if strings.ToLower(contenido.GetGenero()) == genero {
			resultado = append(resultado, contenido)
		}
	}
	return resultado
}

func BuscarPorID(catalogo []models.Reproducible, id int) (models.Reproducible, error) {
	for _, contenido := range catalogo {
		if contenido.GetID() == id {
			return contenido, nil
		}
	}
	return nil, errors.New("contenido no encontrado")
}

// GenerosDisponibles devuelve la lista de géneros únicos presentes en el
// catálogo, usada por el servicio web GET /categorias.
func GenerosDisponibles(catalogo []models.Reproducible) []string {
	vistos := make(map[string]bool)
	generos := make([]string, 0)

	for _, contenido := range catalogo {
		genero := contenido.GetGenero()
		if !vistos[genero] {
			vistos[genero] = true
			generos = append(generos, genero)
		}
	}
	return generos
}
