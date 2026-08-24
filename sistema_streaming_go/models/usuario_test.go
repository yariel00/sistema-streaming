package models

import "testing"

// --- Pruebas unitarias ---------------------------------------------------

func TestNuevoUsuario_DatosValidos(t *testing.T) {
	u, err := NuevoUsuario(1, "Ana Pérez", "ana@correo.com")
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if u.GetID() != 1 || u.GetNombre() != "Ana Pérez" || u.GetEmail() != "ana@correo.com" {
		t.Errorf("los datos del usuario no coinciden con los ingresados")
	}
	if u.GetSuscripcion() != "Sin suscripción" {
		t.Errorf("un usuario nuevo debería iniciar sin suscripción, obtuvo: %s", u.GetSuscripcion())
	}
}

func TestNuevoUsuario_IDInvalido(t *testing.T) {
	_, err := NuevoUsuario(0, "Ana", "ana@correo.com")
	if err == nil {
		t.Error("se esperaba un error para un ID menor o igual a cero")
	}
}

func TestNuevoUsuario_NombreVacio(t *testing.T) {
	_, err := NuevoUsuario(1, "   ", "ana@correo.com")
	if err == nil {
		t.Error("se esperaba un error para un nombre vacío")
	}
}

func TestNuevoUsuario_EmailInvalido(t *testing.T) {
	_, err := NuevoUsuario(1, "Ana", "ana-correo-sin-arroba")
	if err == nil {
		t.Error("se esperaba un error para un correo sin '@'")
	}
}

func TestUsuario_SetSuscripcion(t *testing.T) {
	u, _ := NuevoUsuario(1, "Ana", "ana@correo.com")

	if err := u.SetSuscripcion("Premium"); err != nil {
		t.Fatalf("no se esperaba error al asignar un plan válido: %v", err)
	}
	if !u.TieneSuscripcion() {
		t.Error("TieneSuscripcion() debería ser true después de asignar un plan")
	}

	if err := u.SetSuscripcion("   "); err == nil {
		t.Error("se esperaba un error al asignar un plan vacío")
	}
}
