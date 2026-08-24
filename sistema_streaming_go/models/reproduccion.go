package models

import "time"

// Reproduccion representa un evento de reproducción registrado en el
// historial de la plataforma. Se usa tanto para el reporte de
// estadísticas como para el servicio web GET /reproducciones.
type Reproduccion struct {
	IDUsuario   int       `json:"id_usuario"`
	IDContenido int       `json:"id_contenido"`
	Titulo      string    `json:"titulo"`
	Tipo        string    `json:"tipo"`
	Fecha       time.Time `json:"fecha"`
}
