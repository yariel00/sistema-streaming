package services

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"sistema-streaming/models"
)

// -----------------------------------------------------------------------
// Este archivo implementa los servicios web (REST) de la Unidad 4.
// Cada endpoint responde con datos serializados en formato JSON mediante
// el paquete estándar encoding/json.
//
// Concurrencia: net/http atiende cada petición entrante en su propia
// goroutine automáticamente. Eso significa que, si dos clientes llaman
// a /suscripciones o /reproducciones al mismo tiempo, ambas peticiones
// se ejecutan en paralelo. Para que eso sea seguro, todo el estado
// compartido (usuarios, historial) vive dentro de PlataformaStreaming y
// está protegido por un sync.RWMutex (ver services/streaming.go), en
// lugar de usar variables globales sin protección.
// -----------------------------------------------------------------------

type respuestaError struct {
	Error string `json:"error"`
}

func escribirJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("error al serializar JSON:", err)
	}
}

func escribirError(w http.ResponseWriter, status int, mensaje string) {
	escribirJSON(w, status, respuestaError{Error: mensaje})
}

// contenidosADTO convierte un slice de Reproducible a un slice de DTOs
// serializables.
func contenidosADTO(catalogo []models.Reproducible) []models.ContenidoDTO {
	dtos := make([]models.ContenidoDTO, 0, len(catalogo))
	for _, c := range catalogo {
		dtos = append(dtos, models.NuevoContenidoDTO(c))
	}
	return dtos
}

func usuariosADTO(usuarios []*models.Usuario) []models.UsuarioDTO {
	dtos := make([]models.UsuarioDTO, 0, len(usuarios))
	for _, u := range usuarios {
		dtos = append(dtos, models.NuevoUsuarioDTO(u))
	}
	return dtos
}

// ------------------------ 1. /usuarios ---------------------------------

type crearUsuarioRequest struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

func (p *PlataformaStreaming) handleUsuarios(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		escribirJSON(w, http.StatusOK, usuariosADTO(p.GetUsuarios()))

	case http.MethodPost:
		var req crearUsuarioRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			escribirError(w, http.StatusBadRequest, "cuerpo de la petición inválido")
			return
		}

		usuario, err := models.NuevoUsuario(req.ID, req.Nombre, req.Email)
		if err != nil {
			escribirError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := p.RegistrarUsuario(usuario); err != nil {
			escribirError(w, http.StatusConflict, err.Error())
			return
		}

		escribirJSON(w, http.StatusCreated, models.NuevoUsuarioDTO(usuario))

	default:
		escribirError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ------------------------ 2. /contenidos --------------------------------

// GET /contenidos                 -> catálogo completo
// GET /contenidos?titulo=coco     -> búsqueda por título
// GET /contenidos?genero=Drama    -> filtro por género
func (p *PlataformaStreaming) handleContenidos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		escribirError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	catalogo := p.GetCatalogo()

	if titulo := r.URL.Query().Get("titulo"); titulo != "" {
		catalogo = BuscarPorTitulo(catalogo, titulo)
	} else if genero := r.URL.Query().Get("genero"); genero != "" {
		catalogo = FiltrarPorGenero(catalogo, genero)
	}

	escribirJSON(w, http.StatusOK, contenidosADTO(catalogo))
}

// ------------------------ 3. /suscripciones ------------------------------

type suscripcionRequest struct {
	IDUsuario int    `json:"id_usuario"`
	Plan      string `json:"plan"`
}

func (p *PlataformaStreaming) handleSuscripciones(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		escribirJSON(w, http.StatusOK, usuariosADTO(p.GetUsuariosSuscritos()))

	case http.MethodPost:
		var req suscripcionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			escribirError(w, http.StatusBadRequest, "cuerpo de la petición inválido")
			return
		}

		if err := p.SuscribirUsuario(req.IDUsuario, req.Plan); err != nil {
			escribirError(w, http.StatusBadRequest, err.Error())
			return
		}

		usuario, _ := p.ObtenerUsuario(req.IDUsuario)
		escribirJSON(w, http.StatusOK, models.NuevoUsuarioDTO(usuario))

	default:
		escribirError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ------------------------ 4. /reproducciones ------------------------------

type reproduccionRequest struct {
	IDUsuario   int `json:"id_usuario"`
	IDContenido int `json:"id_contenido"`
}

func (p *PlataformaStreaming) handleReproducciones(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		escribirJSON(w, http.StatusOK, p.GetHistorial())

	case http.MethodPost:
		var req reproduccionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			escribirError(w, http.StatusBadRequest, "cuerpo de la petición inválido")
			return
		}

		mensaje, err := p.ReproducirContenido(req.IDUsuario, req.IDContenido)
		if err != nil {
			escribirError(w, http.StatusBadRequest, err.Error())
			return
		}

		escribirJSON(w, http.StatusOK, map[string]string{"mensaje": mensaje})

	default:
		escribirError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ------------------------ 5. /reportes ------------------------------------

func (p *PlataformaStreaming) handleReportes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		escribirError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}
	escribirJSON(w, http.StatusOK, p.GenerarReporte())
}

// ------------------------ 6. /estado ----------------------------------------

func (p *PlataformaStreaming) handleEstado(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		escribirError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}
	escribirJSON(w, http.StatusOK, map[string]interface{}{
		"estado":     "activo",
		"plataforma": p.GetNombre(),
		"hora_servidor": time.Now().Format(time.RFC3339),
	})
}

// ------------------------ 7. /planes ----------------------------------------

func (p *PlataformaStreaming) handlePlanes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		escribirError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}
	escribirJSON(w, http.StatusOK, PlanesDisponibles())
}

// ------------------------ 8. /categorias ------------------------------------

func (p *PlataformaStreaming) handleCategorias(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		escribirError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}
	escribirJSON(w, http.StatusOK, GenerosDisponibles(p.GetCatalogo()))
}

// IniciarServidorWeb registra las rutas y arranca el servidor HTTP.
// Se ejecuta en su propia goroutine desde main.go para no bloquear el
// menú de consola: la aplicación puede usarse por consola y por HTTP
// al mismo tiempo, ambos accediendo de forma segura al mismo estado
// gracias al mutex de PlataformaStreaming.
func (p *PlataformaStreaming) IniciarServidorWeb(puerto string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/usuarios", p.handleUsuarios)
	mux.HandleFunc("/contenidos", p.handleContenidos)
	mux.HandleFunc("/suscripciones", p.handleSuscripciones)
	mux.HandleFunc("/reproducciones", p.handleReproducciones)
	mux.HandleFunc("/reportes", p.handleReportes)
	mux.HandleFunc("/estado", p.handleEstado)
	mux.HandleFunc("/planes", p.handlePlanes)
	mux.HandleFunc("/categorias", p.handleCategorias)

	log.Printf("Servidor web escuchando en http://localhost%s\n", puerto)
	return http.ListenAndServe(puerto, mux)
}
