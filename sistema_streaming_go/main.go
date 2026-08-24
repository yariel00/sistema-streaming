package main

import (
	"fmt"
	"sync"

	"sistema-streaming/models"
	"sistema-streaming/services"
	"sistema-streaming/utils"
)

func mostrarCatalogo(catalogo []models.Reproducible) {
	if len(catalogo) == 0 {
		fmt.Println("No se encontraron contenidos.")
		return
	}

	fmt.Println("\n--- CATÁLOGO ---")
	for _, contenido := range catalogo {
		fmt.Printf(
			"ID: %d | %s | %s | Género: %s\n",
			contenido.GetID(),
			contenido.GetTipo(),
			contenido.GetTitulo(),
			contenido.GetGenero(),
		)
	}
}

// simularReproduccionesConcurrentes lanza varias reproducciones al mismo
// tiempo usando goroutines, para demostrar que PlataformaStreaming
// soporta acceso concurrente de forma segura (gracias al sync.RWMutex
// interno). Se usa un sync.WaitGroup para esperar a que todas terminen
// antes de continuar.
func simularReproduccionesConcurrentes(plataforma *services.PlataformaStreaming) {
	solicitudes := []struct {
		idUsuario   int
		idContenido int
	}{
		{1, 1}, {1, 4}, {2, 2}, {2, 5}, {3, 3}, {3, 6},
	}

	var wg sync.WaitGroup
	var mu sync.Mutex // protege únicamente la impresión ordenada en consola

	fmt.Println("\n--- Simulando reproducciones concurrentes ---")
	for i, s := range solicitudes {
		wg.Add(1)
		go func(indice, idUsuario, idContenido int) {
			defer wg.Done()

			mensaje, err := plataforma.ReproducirContenido(idUsuario, idContenido)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				fmt.Printf("[goroutine %d] Error (usuario %d, contenido %d): %v\n",
					indice, idUsuario, idContenido, err)
				return
			}
			fmt.Printf("[goroutine %d] %s\n", indice, mensaje)
		}(i, s.idUsuario, s.idContenido)
	}

	wg.Wait()
	fmt.Println("--- Simulación finalizada: todas las goroutines terminaron ---")
	fmt.Println("Nota: los usuarios 1, 2 y 3 deben existir y tener una suscripción")
	fmt.Println("activa (opciones 4 y 5 del menú) para que la simulación no falle.")
}

func main() {
	catalogoInicial := []models.Reproducible{
		models.Pelicula{ID: 1, Titulo: "Interstellar", Genero: "Ciencia ficción", Duracion: 169},
		models.Pelicula{ID: 2, Titulo: "Coco", Genero: "Animación", Duracion: 105},
		models.Pelicula{ID: 3, Titulo: "El Conjuro", Genero: "Terror", Duracion: 112},
		models.Serie{ID: 4, Titulo: "Stranger Things", Genero: "Ciencia ficción", Temporadas: 4},
		models.Serie{ID: 5, Titulo: "Breaking Bad", Genero: "Drama", Temporadas: 5},
		models.Serie{ID: 6, Titulo: "Dark", Genero: "Ciencia ficción", Temporadas: 3},
	}

	plataforma := services.NuevaPlataforma("GoStream", catalogoInicial)
	servidorIniciado := false

	for {
		fmt.Printf("\n=== %s ===\n", plataforma.GetNombre())
		fmt.Println("1. Mostrar catálogo")
		fmt.Println("2. Buscar contenido por título")
		fmt.Println("3. Filtrar por género")
		fmt.Println("4. Registrar usuario")
		fmt.Println("5. Activar suscripción")
		fmt.Println("6. Reproducir contenido")
		fmt.Println("7. Ver datos de usuario")
		fmt.Println("8. Iniciar servidor web (servicios REST en :8080)")
		fmt.Println("9. Simular reproducciones concurrentes (demo)")
		fmt.Println("0. Salir")

		opcion, err := utils.LeerEntero("Seleccione una opción: ")
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		switch opcion {
		case 1:
			mostrarCatalogo(plataforma.GetCatalogo())

		case 2:
			texto := utils.LeerTexto("Ingrese el título o parte del título: ")
			resultados := services.BuscarPorTitulo(plataforma.GetCatalogo(), texto)
			mostrarCatalogo(resultados)

		case 3:
			genero := utils.LeerTexto("Ingrese el género: ")
			resultados := services.FiltrarPorGenero(plataforma.GetCatalogo(), genero)
			mostrarCatalogo(resultados)

		case 4:
			id, err := utils.LeerEntero("ID del usuario: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			nombre := utils.LeerTexto("Nombre: ")
			email := utils.LeerTexto("Correo: ")

			usuario, err := models.NuevoUsuario(id, nombre, email)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			if err := plataforma.RegistrarUsuario(usuario); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println("Usuario registrado correctamente.")

		case 5:
			id, err := utils.LeerEntero("ID del usuario: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			plan := utils.LeerTexto("Plan (Básico / Estándar / Premium): ")

			if err := plataforma.SuscribirUsuario(id, plan); err != nil {
				fmt.Println("Error:", err)
				continue
			}

		case 6:
			idUsuario, err := utils.LeerEntero("ID del usuario: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			idContenido, err := utils.LeerEntero("ID del contenido: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			mensaje, err := plataforma.ReproducirContenido(idUsuario, idContenido)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println(mensaje)

		case 7:
			id, err := utils.LeerEntero("ID del usuario: ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			usuario, err := plataforma.ObtenerUsuario(id)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Printf(
				"ID: %d | Nombre: %s | Correo: %s | Suscripción: %s\n",
				usuario.GetID(),
				usuario.GetNombre(),
				usuario.GetEmail(),
				usuario.GetSuscripcion(),
			)

		case 8:
			if servidorIniciado {
				fmt.Println("El servidor web ya está en ejecución en http://localhost:8080")
				continue
			}
			servidorIniciado = true
			go func() {
				if err := plataforma.IniciarServidorWeb(":8080"); err != nil {
					fmt.Println("Error al iniciar el servidor web:", err)
				}
			}()
			fmt.Println("Servidor web iniciado en http://localhost:8080")
			fmt.Println("Prueba, por ejemplo: http://localhost:8080/contenidos")

		case 9:
			simularReproduccionesConcurrentes(plataforma)

		case 0:
			fmt.Println("Programa finalizado.")
			return

		default:
			fmt.Println("Opción no válida.")
		}
	}
}
