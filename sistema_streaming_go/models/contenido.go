package models

import "fmt"

// Reproducible define el comportamiento común de cualquier contenido
// que pueda ser mostrado y reproducido dentro de la plataforma.
type Reproducible interface {
	GetID() int
	GetTitulo() string
	GetGenero() string
	GetTipo() string
	Reproducir() string
}

// Pelicula representa una película dentro del catálogo.
type Pelicula struct {
	ID       int
	Titulo   string
	Genero   string
	Duracion int // minutos
}

func (p Pelicula) GetID() int        { return p.ID }
func (p Pelicula) GetTitulo() string { return p.Titulo }
func (p Pelicula) GetGenero() string { return p.Genero }
func (p Pelicula) GetTipo() string   { return "Película" }

func (p Pelicula) Reproducir() string {
	return fmt.Sprintf("Reproduciendo película: %s (%d minutos)", p.Titulo, p.Duracion)
}

// Serie representa una serie dentro del catálogo.
type Serie struct {
	ID         int
	Titulo     string
	Genero     string
	Temporadas int
}

func (s Serie) GetID() int        { return s.ID }
func (s Serie) GetTitulo() string { return s.Titulo }
func (s Serie) GetGenero() string { return s.Genero }
func (s Serie) GetTipo() string   { return "Serie" }

func (s Serie) Reproducir() string {
	return fmt.Sprintf("Reproduciendo serie: %s (%d temporadas)", s.Titulo, s.Temporadas)
}

// ContenidoDTO es la representación serializable en JSON de un contenido
// del catálogo. Los structs Pelicula y Serie no tienen campos exportados
// homogéneos entre sí (Duracion vs Temporadas), por lo que se usa un DTO
// (Data Transfer Object) para exponerlos de forma uniforme en los
// servicios web.
type ContenidoDTO struct {
	ID     int    `json:"id"`
	Tipo   string `json:"tipo"`
	Titulo string `json:"titulo"`
	Genero string `json:"genero"`
}

// NuevoContenidoDTO construye un DTO a partir de cualquier valor que
// implemente la interfaz Reproducible.
func NuevoContenidoDTO(c Reproducible) ContenidoDTO {
	return ContenidoDTO{
		ID:     c.GetID(),
		Tipo:   c.GetTipo(),
		Titulo: c.GetTitulo(),
		Genero: c.GetGenero(),
	}
}
