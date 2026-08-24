package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sistema-streaming/models"
)

// --- Pruebas de aceptación --------------------------------------------------
// Simulan el uso real de los servicios web tal como lo haría un cliente
// HTTP (por ejemplo Postman o un navegador), usando httptest para no
// necesitar un servidor real escuchando en un puerto.

func plataformaWebPrueba() *PlataformaStreaming {
	catalogo := []models.Reproducible{
		models.Pelicula{ID: 1, Titulo: "Interstellar", Genero: "Ciencia ficción", Duracion: 169},
		models.Serie{ID: 2, Titulo: "Dark", Genero: "Ciencia ficción", Temporadas: 3},
	}
	return NuevaPlataforma("GoStream Test", catalogo)
}

func TestServicioWeb_GetContenidos(t *testing.T) {
	p := plataformaWebPrueba()

	req := httptest.NewRequest(http.MethodGet, "/contenidos", nil)
	rec := httptest.NewRecorder()

	p.handleContenidos(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba status 200, se obtuvo %d", rec.Code)
	}

	var resultado []models.ContenidoDTO
	if err := json.NewDecoder(rec.Body).Decode(&resultado); err != nil {
		t.Fatalf("la respuesta no es un JSON válido: %v", err)
	}
	if len(resultado) != 2 {
		t.Errorf("se esperaban 2 contenidos, se obtuvo %d", len(resultado))
	}
}

func TestServicioWeb_GetContenidos_FiltroPorGenero(t *testing.T) {
	p := plataformaWebPrueba()

	req := httptest.NewRequest(http.MethodGet, "/contenidos?genero=Ciencia%20ficci%C3%B3n", nil)
	rec := httptest.NewRecorder()

	p.handleContenidos(rec, req)

	var resultado []models.ContenidoDTO
	_ = json.NewDecoder(rec.Body).Decode(&resultado)
	if len(resultado) != 2 {
		t.Errorf("se esperaban 2 contenidos de 'Ciencia ficción', se obtuvo %d", len(resultado))
	}
}

func TestServicioWeb_PostUsuario_CreaUsuario(t *testing.T) {
	p := plataformaWebPrueba()

	cuerpo, _ := json.Marshal(crearUsuarioRequest{ID: 1, Nombre: "Ana", Email: "ana@correo.com"})
	req := httptest.NewRequest(http.MethodPost, "/usuarios", bytes.NewReader(cuerpo))
	rec := httptest.NewRecorder()

	p.handleUsuarios(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba status 201, se obtuvo %d. Cuerpo: %s", rec.Code, rec.Body.String())
	}

	var dto models.UsuarioDTO
	if err := json.NewDecoder(rec.Body).Decode(&dto); err != nil {
		t.Fatalf("la respuesta no es un JSON válido: %v", err)
	}
	if dto.Nombre != "Ana" {
		t.Errorf("se esperaba el usuario 'Ana', se obtuvo '%s'", dto.Nombre)
	}
}

func TestServicioWeb_PostUsuario_DatosInvalidos(t *testing.T) {
	p := plataformaWebPrueba()

	cuerpo, _ := json.Marshal(crearUsuarioRequest{ID: 0, Nombre: "", Email: "correo-invalido"})
	req := httptest.NewRequest(http.MethodPost, "/usuarios", bytes.NewReader(cuerpo))
	rec := httptest.NewRecorder()

	p.handleUsuarios(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("se esperaba status 400 para datos inválidos, se obtuvo %d", rec.Code)
	}
}

func TestServicioWeb_FlujoCompleto_SuscripcionYReproduccion(t *testing.T) {
	p := plataformaWebPrueba()

	// 1. Crear usuario
	cuerpoUsuario, _ := json.Marshal(crearUsuarioRequest{ID: 1, Nombre: "Ana", Email: "ana@correo.com"})
	reqUsuario := httptest.NewRequest(http.MethodPost, "/usuarios", bytes.NewReader(cuerpoUsuario))
	recUsuario := httptest.NewRecorder()
	p.handleUsuarios(recUsuario, reqUsuario)
	if recUsuario.Code != http.StatusCreated {
		t.Fatalf("no se pudo crear el usuario: %s", recUsuario.Body.String())
	}

	// 2. Suscribirlo
	cuerpoSus, _ := json.Marshal(suscripcionRequest{IDUsuario: 1, Plan: "Premium"})
	reqSus := httptest.NewRequest(http.MethodPost, "/suscripciones", bytes.NewReader(cuerpoSus))
	recSus := httptest.NewRecorder()
	p.handleSuscripciones(recSus, reqSus)
	if recSus.Code != http.StatusOK {
		t.Fatalf("no se pudo suscribir al usuario: %s", recSus.Body.String())
	}

	// 3. Reproducir contenido
	cuerpoRepro, _ := json.Marshal(reproduccionRequest{IDUsuario: 1, IDContenido: 1})
	reqRepro := httptest.NewRequest(http.MethodPost, "/reproducciones", bytes.NewReader(cuerpoRepro))
	recRepro := httptest.NewRecorder()
	p.handleReproducciones(recRepro, reqRepro)
	if recRepro.Code != http.StatusOK {
		t.Fatalf("no se pudo reproducir el contenido: %s", recRepro.Body.String())
	}

	// 4. Verificar el reporte final
	reqReporte := httptest.NewRequest(http.MethodGet, "/reportes", nil)
	recReporte := httptest.NewRecorder()
	p.handleReportes(recReporte, reqReporte)

	var reporte ReporteGeneral
	_ = json.NewDecoder(recReporte.Body).Decode(&reporte)
	if reporte.TotalUsuarios != 1 || reporte.UsuariosConSuscripcion != 1 || reporte.TotalReproducciones != 1 {
		t.Errorf("el reporte final no refleja el flujo esperado: %+v", reporte)
	}
}

func TestServicioWeb_MetodoNoPermitido(t *testing.T) {
	p := plataformaWebPrueba()

	req := httptest.NewRequest(http.MethodDelete, "/contenidos", nil)
	rec := httptest.NewRecorder()

	p.handleContenidos(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("se esperaba status 405, se obtuvo %d", rec.Code)
	}
}
