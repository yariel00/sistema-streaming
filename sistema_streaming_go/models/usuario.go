package models

import (
	"errors"
	"strings"
)

// Usuario aplica encapsulación: sus atributos son privados y se accede
// a ellos mediante métodos públicos.
type Usuario struct {
	id          int
	nombre      string
	email       string
	suscripcion string
}

// NuevoUsuario funciona como constructor y valida los datos antes de crear
// el objeto.
func NuevoUsuario(id int, nombre, email string) (*Usuario, error) {
	nombre = strings.TrimSpace(nombre)
	email = strings.TrimSpace(email)

	if id <= 0 {
		return nil, errors.New("el ID del usuario debe ser mayor que cero")
	}
	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("el correo electrónico no es válido")
	}

	return &Usuario{
		id:          id,
		nombre:      nombre,
		email:       email,
		suscripcion: "Sin suscripción",
	}, nil
}

func (u *Usuario) GetID() int             { return u.id }
func (u *Usuario) GetNombre() string      { return u.nombre }
func (u *Usuario) GetEmail() string       { return u.email }
func (u *Usuario) GetSuscripcion() string { return u.suscripcion }
func (u *Usuario) TieneSuscripcion() bool { return u.suscripcion != "Sin suscripción" }
func (u *Usuario) SetSuscripcion(plan string) error {
	plan = strings.TrimSpace(plan)
	if plan == "" {
		return errors.New("el plan no puede estar vacío")
	}
	u.suscripcion = plan
	return nil
}

// UsuarioDTO es la representación serializable en JSON de un Usuario.
// Como los campos de Usuario son privados (encapsulación), no se pueden
// serializar directamente con encoding/json: se necesita este DTO con
// campos exportados para exponer el usuario en los servicios web.
type UsuarioDTO struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Email       string `json:"email"`
	Suscripcion string `json:"suscripcion"`
}

// NuevoUsuarioDTO construye un DTO a partir de un *Usuario.
func NuevoUsuarioDTO(u *Usuario) UsuarioDTO {
	return UsuarioDTO{
		ID:          u.GetID(),
		Nombre:      u.GetNombre(),
		Email:       u.GetEmail(),
		Suscripcion: u.GetSuscripcion(),
	}
}
