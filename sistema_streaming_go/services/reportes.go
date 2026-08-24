package services

// ReporteGeneral agrupa estadísticas generales de la plataforma.
// Se serializa como JSON en el servicio web GET /reportes.
type ReporteGeneral struct {
	NombrePlataforma       string   `json:"nombre_plataforma"`
	TotalUsuarios          int      `json:"total_usuarios"`
	UsuariosConSuscripcion int      `json:"usuarios_con_suscripcion"`
	TotalContenidos        int      `json:"total_contenidos"`
	TotalReproducciones    int      `json:"total_reproducciones"`
	GenerosDisponibles     []string `json:"generos_disponibles"`
}

// GenerarReporte recopila estadísticas de todos los módulos de la
// plataforma (usuarios, catálogo, suscripciones y reproducciones).
func (p *PlataformaStreaming) GenerarReporte() ReporteGeneral {
	catalogo := p.GetCatalogo()
	usuarios := p.GetUsuarios()
	suscritos := p.GetUsuariosSuscritos()
	historial := p.GetHistorial()

	return ReporteGeneral{
		NombrePlataforma:       p.GetNombre(),
		TotalUsuarios:          len(usuarios),
		UsuariosConSuscripcion: len(suscritos),
		TotalContenidos:        len(catalogo),
		TotalReproducciones:    len(historial),
		GenerosDisponibles:     GenerosDisponibles(catalogo),
	}
}
