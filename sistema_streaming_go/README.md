# Sistema de Gestión de Streaming — GoStream

**Materia:** Programación Orientada a Objetos
**Carrera:** [COMPLETAR: Carrera]
**Institución:** [COMPLETAR: Institución]
**Integrantes del grupo:**
- [COMPLETAR: Nombre completo 1]
- [COMPLETAR: Nombre completo 2 (si aplica)]

**Fecha:** [COMPLETAR: dd/mm/aaaa]

---

## Objetivo del programa

Desarrollar un sistema de gestión de streaming en Go (Golang) que permita administrar
usuarios, un catálogo de contenidos (películas y series), suscripciones, reproducciones
y reportes generales, exponiendo esta funcionalidad tanto por consola como mediante
**servicios web REST con serialización JSON**, aplicando programación orientada a
objetos, manejo de errores, interfaces y concurrencia segura.

## Funcionalidades principales

### Por consola (menú interactivo)
1. Mostrar catálogo de contenidos.
2. Buscar contenido por título.
3. Filtrar contenido por género.
4. Registrar usuario.
5. Activar suscripción (Básico / Estándar / Premium).
6. Reproducir contenido (requiere suscripción activa).
7. Ver datos de un usuario.
8. Iniciar el servidor web de servicios REST (puerto `:8080`).
9. Simular reproducciones concurrentes (demo de concurrencia con goroutines).
0. Salir.

### Servicios web REST (JSON)

| # | Endpoint          | Método | Descripción                                             |
|---|--------------------|--------|-----------------------------------------------------------|
| 1 | `/usuarios`        | GET    | Lista todos los usuarios registrados                      |
| 1 | `/usuarios`        | POST   | Registra un nuevo usuario                                  |
| 2 | `/contenidos`      | GET    | Lista el catálogo (admite `?titulo=` y `?genero=`)         |
| 3 | `/suscripciones`   | GET    | Lista los usuarios con suscripción activa                  |
| 3 | `/suscripciones`   | POST   | Activa la suscripción de un usuario                        |
| 4 | `/reproducciones`  | GET    | Lista el historial de reproducciones                       |
| 4 | `/reproducciones`  | POST   | Registra la reproducción de un contenido                   |
| 5 | `/reportes`        | GET    | Estadísticas generales de la plataforma                    |
| 6 | `/estado`          | GET    | Estado actual del sistema                                  |
| 7 | `/planes`          | GET    | Planes de suscripción disponibles y sus precios             |
| 8 | `/categorias`      | GET    | Géneros disponibles en el catálogo                          |

Todas las respuestas se serializan en formato **JSON** mediante el paquete estándar
`encoding/json`.

## Estructura del proyecto

```
sistema_streaming_go/
├── main.go                     # Menú de consola y arranque del servidor web
├── go.mod
├── models/
│   ├── contenido.go            # Interfaz Reproducible, Pelicula, Serie, ContenidoDTO
│   ├── usuario.go              # Usuario (encapsulado), UsuarioDTO
│   └── reproduccion.go         # Registro de historial de reproducciones
├── services/
│   ├── catalogo.go             # Búsqueda, filtrado y utilidades del catálogo
│   ├── streaming.go            # PlataformaStreaming (lógica central + mutex)
│   ├── reportes.go             # Generación de reportes estadísticos
│   ├── webserver.go            # Servicios web REST (Unidad 4)
│   ├── catalogo_test.go        # Pruebas unitarias
│   ├── streaming_test.go       # Pruebas de integración y concurrencia
│   └── webserver_test.go       # Pruebas de aceptación (HTTP end-to-end)
├── models/usuario_test.go      # Pruebas unitarias de Usuario
└── utils/
    └── input.go                # Lectura de entrada por consola
```

## Cómo ejecutar

```bash
go run .
```

Desde el menú, selecciona la opción **8** para levantar el servidor web en
`http://localhost:8080` y probar los endpoints (por ejemplo, con el navegador,
Postman o `curl`):

```bash
curl http://localhost:8080/contenidos
curl http://localhost:8080/planes
curl -X POST http://localhost:8080/usuarios \
  -H "Content-Type: application/json" \
  -d '{"id":1,"nombre":"Ana Pérez","email":"ana@correo.com"}'
```

## Cómo ejecutar las pruebas

```bash
go test ./...
go test -race ./...   # verifica además la ausencia de condiciones de carrera
```

## Concurrencia

El servidor HTTP (`net/http`) atiende cada petición en su propia goroutine.
`PlataformaStreaming` protege su estado compartido (usuarios e historial de
reproducciones) con un `sync.RWMutex`, por lo que múltiples usuarios pueden
registrarse, suscribirse o reproducir contenido de forma simultánea sin
riesgo de condiciones de carrera. La opción **9** del menú y la prueba
`TestAccesoConcurrente_SinCondicionesDeCarrera` (`services/streaming_test.go`)
demuestran este comportamiento con múltiples goroutines concurrentes.

## Integración de las 4 unidades

- **Unidad 1:** funciones, paquetes y estructura modular del proyecto.
- **Unidad 2:** structs, métodos, constructores, slices y maps.
- **Unidad 3:** encapsulación, interfaces (`Reproducible`) y manejo de errores.
- **Unidad 4:** servicios web REST, `net/http`, serialización JSON y concurrencia.

## Repositorio

`[COMPLETAR: enlace al repositorio de GitHub]`
