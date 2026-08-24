package services

import (
	"testing"

	"sistema-streaming/models"
)

func catalogoPrueba() []models.Reproducible {
	return []models.Reproducible{
		models.Pelicula{ID: 1, Titulo: "Interstellar", Genero: "Ciencia ficción", Duracion: 169},
		models.Pelicula{ID: 2, Titulo: "Coco", Genero: "Animación", Duracion: 105},
		models.Serie{ID: 3, Titulo: "Dark", Genero: "Ciencia ficción", Temporadas: 3},
	}
}

// --- Pruebas unitarias ---------------------------------------------------

func TestBuscarPorTitulo(t *testing.T) {
	resultados := BuscarPorTitulo(catalogoPrueba(), "co")
	if len(resultados) != 1 || resultados[0].GetTitulo() != "Coco" {
		t.Errorf("se esperaba encontrar solo 'Coco', se obtuvo: %v", resultados)
	}
}

func TestBuscarPorTitulo_SinCoincidencias(t *testing.T) {
	resultados := BuscarPorTitulo(catalogoPrueba(), "inexistente")
	if len(resultados) != 0 {
		t.Errorf("se esperaba una lista vacía, se obtuvo %d resultados", len(resultados))
	}
}

func TestFiltrarPorGenero(t *testing.T) {
	resultados := FiltrarPorGenero(catalogoPrueba(), "Ciencia ficción")
	if len(resultados) != 2 {
		t.Errorf("se esperaban 2 resultados de 'Ciencia ficción', se obtuvo %d", len(resultados))
	}
}

func TestBuscarPorID_Encontrado(t *testing.T) {
	contenido, err := BuscarPorID(catalogoPrueba(), 2)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if contenido.GetTitulo() != "Coco" {
		t.Errorf("se esperaba 'Coco', se obtuvo '%s'", contenido.GetTitulo())
	}
}

func TestBuscarPorID_NoEncontrado(t *testing.T) {
	_, err := BuscarPorID(catalogoPrueba(), 999)
	if err == nil {
		t.Error("se esperaba un error para un ID inexistente")
	}
}

func TestGenerosDisponibles(t *testing.T) {
	generos := GenerosDisponibles(catalogoPrueba())
	if len(generos) != 2 {
		t.Errorf("se esperaban 2 géneros únicos, se obtuvo %d: %v", len(generos), generos)
	}
}

func TestPrecioPlan_PlanesValidos(t *testing.T) {
	casos := map[string]float64{
		"Básico":   5.99,
		"Estándar": 8.99,
		"Premium":  12.99,
	}
	for plan, esperado := range casos {
		precio, err := PrecioPlan(plan)
		if err != nil {
			t.Fatalf("no se esperaba error para el plan %s: %v", plan, err)
		}
		if precio != esperado {
			t.Errorf("plan %s: se esperaba %.2f, se obtuvo %.2f", plan, esperado, precio)
		}
	}
}

func TestPrecioPlan_PlanInvalido(t *testing.T) {
	_, err := PrecioPlan("Gold")
	if err == nil {
		t.Error("se esperaba un error para un plan no reconocido")
	}
}
