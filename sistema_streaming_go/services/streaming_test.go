package services

import (
	"sync"
	"testing"

	"sistema-streaming/models"
)

func plataformaPrueba() *PlataformaStreaming {
	catalogo := []models.Reproducible{
		models.Pelicula{ID: 1, Titulo: "Interstellar", Genero: "Ciencia ficción", Duracion: 169},
		models.Serie{ID: 2, Titulo: "Dark", Genero: "Ciencia ficción", Temporadas: 3},
	}
	return NuevaPlataforma("GoStream Test", catalogo)
}

// --- Pruebas de integración -----------------------------------------------
// Estas pruebas verifican que varios componentes (models.Usuario,
// PlataformaStreaming, PrecioPlan) trabajen correctamente en conjunto.

func TestFlujoCompleto_RegistroSuscripcionYReproduccion(t *testing.T) {
	p := plataformaPrueba()

	usuario, err := models.NuevoUsuario(1, "Ana", "ana@correo.com")
	if err != nil {
		t.Fatalf("no se esperaba error al crear el usuario: %v", err)
	}

	if err := p.RegistrarUsuario(usuario); err != nil {
		t.Fatalf("no se esperaba error al registrar el usuario: %v", err)
	}

	if err := p.SuscribirUsuario(1, "Premium"); err != nil {
		t.Fatalf("no se esperaba error al suscribir al usuario: %v", err)
	}

	mensaje, err := p.ReproducirContenido(1, 1)
	if err != nil {
		t.Fatalf("no se esperaba error al reproducir contenido: %v", err)
	}
	if mensaje == "" {
		t.Error("se esperaba un mensaje de reproducción no vacío")
	}

	historial := p.GetHistorial()
	if len(historial) != 1 {
		t.Fatalf("se esperaba 1 registro en el historial, se obtuvo %d", len(historial))
	}
	if historial[0].IDUsuario != 1 || historial[0].IDContenido != 1 {
		t.Errorf("el registro del historial no coincide con la reproducción realizada")
	}

	reporte := p.GenerarReporte()
	if reporte.TotalUsuarios != 1 || reporte.TotalReproducciones != 1 {
		t.Errorf("el reporte no refleja el estado esperado: %+v", reporte)
	}
}

func TestRegistrarUsuario_IDDuplicado(t *testing.T) {
	p := plataformaPrueba()
	u1, _ := models.NuevoUsuario(1, "Ana", "ana@correo.com")
	u2, _ := models.NuevoUsuario(1, "Beto", "beto@correo.com")

	if err := p.RegistrarUsuario(u1); err != nil {
		t.Fatalf("no se esperaba error en el primer registro: %v", err)
	}
	if err := p.RegistrarUsuario(u2); err == nil {
		t.Error("se esperaba un error al registrar un ID de usuario duplicado")
	}
}

func TestReproducirContenido_SinSuscripcion(t *testing.T) {
	p := plataformaPrueba()
	usuario, _ := models.NuevoUsuario(1, "Ana", "ana@correo.com")
	_ = p.RegistrarUsuario(usuario)

	_, err := p.ReproducirContenido(1, 1)
	if err == nil {
		t.Error("se esperaba un error: el usuario no tiene una suscripción activa")
	}
}

func TestReproducirContenido_UsuarioInexistente(t *testing.T) {
	p := plataformaPrueba()
	_, err := p.ReproducirContenido(999, 1)
	if err == nil {
		t.Error("se esperaba un error para un usuario inexistente")
	}
}

// --- Prueba de concurrencia -------------------------------------------------
// Simula múltiples solicitudes de suscripción y reproducción simultáneas
// (como harían varias peticiones HTTP al mismo tiempo) para verificar que
// PlataformaStreaming no sufre condiciones de carrera. Ejecutar con
// `go test -race ./...` para que el detector de carreras de Go confirme
// la ausencia de accesos concurrentes inseguros.
func TestAccesoConcurrente_SinCondicionesDeCarrera(t *testing.T) {
	p := plataformaPrueba()

	const totalUsuarios = 20
	var wg sync.WaitGroup

	// Registro y suscripción concurrente de varios usuarios.
	for i := 1; i <= totalUsuarios; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			u, err := models.NuevoUsuario(id, "Usuario", "usuario@correo.com")
			if err != nil {
				t.Errorf("error inesperado creando usuario %d: %v", id, err)
				return
			}
			if err := p.RegistrarUsuario(u); err != nil {
				t.Errorf("error inesperado registrando usuario %d: %v", id, err)
				return
			}
			if err := p.SuscribirUsuario(id, "Básico"); err != nil {
				t.Errorf("error inesperado suscribiendo usuario %d: %v", id, err)
			}
		}(i)
	}
	wg.Wait()

	// Reproducciones concurrentes sobre los usuarios ya suscritos.
	for i := 1; i <= totalUsuarios; i++ {
		wg.Add(1)
		go func(idUsuario int) {
			defer wg.Done()
			idContenido := 1
			if idUsuario%2 == 0 {
				idContenido = 2
			}
			if _, err := p.ReproducirContenido(idUsuario, idContenido); err != nil {
				t.Errorf("error inesperado reproduciendo para usuario %d: %v", idUsuario, err)
			}
		}(i)
	}
	wg.Wait()

	if len(p.GetUsuarios()) != totalUsuarios {
		t.Errorf("se esperaban %d usuarios registrados, se obtuvo %d", totalUsuarios, len(p.GetUsuarios()))
	}
	if len(p.GetHistorial()) != totalUsuarios {
		t.Errorf("se esperaban %d reproducciones en el historial, se obtuvo %d", totalUsuarios, len(p.GetHistorial()))
	}
}
