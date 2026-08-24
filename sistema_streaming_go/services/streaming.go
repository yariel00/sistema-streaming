package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"sistema-streaming/models"
)

// PlataformaStreaming concentra la gestión de usuarios y catálogo.
// El campo usuarios es privado para mantener la encapsulación.
//
// mu protege el acceso concurrente a usuarios e historial: el servidor
// web (services/webserver.go) atiende cada petición HTTP en una goroutine
// distinta, por lo que dos usuarios podrían registrarse, suscribirse o
// reproducir contenido al mismo tiempo. Sin este mutex, escrituras
// simultáneas sobre el mapa "usuarios" podrían corromper el estado
// (data race) o hacer que el programa entre en pánico.
type PlataformaStreaming struct {
	mu        sync.RWMutex
	nombre    string
	catalogo  []models.Reproducible
	usuarios  map[int]*models.Usuario
	historial []models.Reproduccion
}

func NuevaPlataforma(nombre string, catalogo []models.Reproducible) *PlataformaStreaming {
	return &PlataformaStreaming{
		nombre:    nombre,
		catalogo:  catalogo,
		usuarios:  make(map[int]*models.Usuario),
		historial: make([]models.Reproduccion, 0),
	}
}

func (p *PlataformaStreaming) GetNombre() string {
	return p.nombre
}

func (p *PlataformaStreaming) GetCatalogo() []models.Reproducible {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Se devuelve una copia para evitar modificaciones externas directas.
	copia := make([]models.Reproducible, len(p.catalogo))
	copy(copia, p.catalogo)
	return copia
}

func (p *PlataformaStreaming) RegistrarUsuario(usuario *models.Usuario) error {
	if usuario == nil {
		return errors.New("el usuario no puede ser nulo")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, existe := p.usuarios[usuario.GetID()]; existe {
		return errors.New("ya existe un usuario con ese ID")
	}
	p.usuarios[usuario.GetID()] = usuario
	return nil
}

func (p *PlataformaStreaming) ObtenerUsuario(id int) (*models.Usuario, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	usuario, existe := p.usuarios[id]
	if !existe {
		return nil, errors.New("usuario no encontrado")
	}
	return usuario, nil
}

// GetUsuarios devuelve una copia de todos los usuarios registrados,
// usada por el servicio web GET /usuarios.
func (p *PlataformaStreaming) GetUsuarios() []*models.Usuario {
	p.mu.RLock()
	defer p.mu.RUnlock()

	lista := make([]*models.Usuario, 0, len(p.usuarios))
	for _, u := range p.usuarios {
		lista = append(lista, u)
	}
	return lista
}

// GetUsuariosSuscritos devuelve solo los usuarios con una suscripción
// activa, usada por el servicio web GET /suscripciones.
func (p *PlataformaStreaming) GetUsuariosSuscritos() []*models.Usuario {
	p.mu.RLock()
	defer p.mu.RUnlock()

	lista := make([]*models.Usuario, 0)
	for _, u := range p.usuarios {
		if u.TieneSuscripcion() {
			lista = append(lista, u)
		}
	}
	return lista
}

func (p *PlataformaStreaming) SuscribirUsuario(id int, plan string) error {
	precio, err := PrecioPlan(plan)
	if err != nil {
		return err
	}

	p.mu.Lock()
	usuario, existe := p.usuarios[id]
	if !existe {
		p.mu.Unlock()
		return errors.New("usuario no encontrado")
	}
	err = usuario.SetSuscripcion(plan)
	p.mu.Unlock()

	if err != nil {
		return err
	}

	fmt.Printf("Suscripción activada. Precio mensual: $%.2f\n", precio)
	return nil
}

func (p *PlataformaStreaming) ReproducirContenido(idUsuario, idContenido int) (string, error) {
	usuario, err := p.ObtenerUsuario(idUsuario)
	if err != nil {
		return "", err
	}

	if !usuario.TieneSuscripcion() {
		return "", errors.New("el usuario necesita una suscripción activa")
	}

	// La lectura del catálogo no requiere el lock general: BuscarPorID
	// solo recorre el slice del catálogo, que nunca se modifica tras la
	// inicialización de la plataforma.
	contenido, err := BuscarPorID(p.catalogo, idContenido)
	if err != nil {
		return "", err
	}

	mensaje := contenido.Reproducir()

	p.mu.Lock()
	p.historial = append(p.historial, models.Reproduccion{
		IDUsuario:   idUsuario,
		IDContenido: idContenido,
		Titulo:      contenido.GetTitulo(),
		Tipo:        contenido.GetTipo(),
		Fecha:       time.Now(),
	})
	p.mu.Unlock()

	return mensaje, nil
}

// GetHistorial devuelve una copia del historial de reproducciones,
// usada por el servicio web GET /reproducciones.
func (p *PlataformaStreaming) GetHistorial() []models.Reproduccion {
	p.mu.RLock()
	defer p.mu.RUnlock()

	copia := make([]models.Reproduccion, len(p.historial))
	copy(copia, p.historial)
	return copia
}

// PrecioPlan es otra función independiente: recibe un plan y devuelve
// un precio, sin modificar el estado de la plataforma.
func PrecioPlan(plan string) (float64, error) {
	switch plan {
	case "Básico", "Basico", "basico", "básico":
		return 5.99, nil
	case "Estándar", "Estandar", "estandar", "estándar":
		return 8.99, nil
	case "Premium", "premium":
		return 12.99, nil
	default:
		return 0, errors.New("plan no válido; use Básico, Estándar o Premium")
	}
}

// PlanesDisponibles devuelve la lista de planes con su precio, usada
// por el servicio web GET /planes.
func PlanesDisponibles() map[string]float64 {
	return map[string]float64{
		"Básico":   5.99,
		"Estándar": 8.99,
		"Premium":  12.99,
	}
}
