---
name: arquitectura
description: Guía para diseñar y desarrollar programas de consola básicos en Go con arquitectura idiomática, separación clara de responsabilidades, manejo explícito de errores y pruebas mantenibles. Usar al crear o refactorizar CLI, comandos, configuración, integración con filesystem o procesos externos.
---

# Skill: arquitectura idiomática de CLI en Go

## Política del proyecto

Las convenciones idiomáticas de Go se consideran requisitos para este proyecto,
no recomendaciones opcionales. Cuando una decisión del producto entre en
conflicto con una convención del lenguaje, documenta la excepción antes de
implementarla. La solución debe preferir la forma más simple, explícita y
comprobable de Go.

## Propósito

Usa esta guía al diseñar, implementar o refactorizar `project-cli`, la CLI que
reemplazará a este proyecto de ejemplo. El objetivo es una herramienta Linux
pequeña, mantenible y predecible para gestionar proyectos locales: listarlos,
crearlos, abrirlos y resolver sus rutas.

La arquitectura recomendada para este proyecto es:

- `cmd/project-cli/main.go` es el entrypoint y convierte un error final en
  código de salida.
- `internal/app` compone configuración, servicios y comandos.
- `internal/command/` contiene fábricas que devuelven `*cobra.Command`.
- Cobra se ocupa del parseo, ayuda y autocompletado; la lógica de negocio no
  debe depender de variables globales de Cobra.

Al crecer el programa, conserva esa composición explícita y separa las
responsabilidades en lugar de introducir un framework o un contenedor de
dependencias.

## Estructura adoptada

Empieza con la estructura mínima y agrega paquetes solo cuando exista una
responsabilidad real:

```text
project-cli/
├── go.mod
├── go.sum
├── cmd/
│   └── project-cli/
│       └── main.go
└── internal/
    ├── app/
    │   └── app.go
    ├── command/
    │   ├── root.go
    │   ├── list.go
    │   ├── open.go
    │   ├── create.go
    │   ├── path.go
    │   └── completion.go
    ├── config/
    │   ├── config.go
    │   └── config_test.go
    ├── services/
    │   ├── service.go
    │   ├── list.go
    │   ├── list_service_test.go
    │   ├── open.go
    │   └── open_service_test.go
    └── process/
        ├── runner.go
        └── runner_test.go
```

Go no exige una división fija entre `cmd` e `internal`. `cmd` es una convención
útil para organizar ejecutables, mientras que `internal` es una restricción del
compilador que evita que otros módulos importen paquetes privados. Este proyecto
adopta ambas convenciones deliberadamente: `cmd/project-cli` contiene solo el
entrypoint y toda la implementación permanece en `internal`. No se crearán
`pkg/`, capas `repository`, `handler` o `usecase` por anticipado.

## Responsabilidades

### `cmd/project-cli/main.go`

El entrypoint debe ser mínimo:

1. Pasa `os.Args[1:]` y los streams a `internal/app`.
2. Escribe el error en `stderr`.
3. Sale con un código distinto de cero si la ejecución falla.

No pongas lógica de filesystem, procesos, configuración ni decisiones de
negocio en el entrypoint. `internal/app` es la raíz de composición y puede
probarse sin terminar el proceso. No uses `log.Fatal` en paquetes reutilizables.

### `internal/app`

`internal/app` carga y valida la configuración, coordina la inicialización
automática cuando corresponde, construye las dependencias concretas y registra
los comandos. No contiene reglas de negocio de proyectos. La composición debe
ser explícita: no se permiten contenedores de dependencias, `init()` con efectos
externos ni estado global mutable.

### `internal/command`

Cada archivo de comando debe exponer una fábrica clara, por ejemplo
`NewListCmd(service services.ProjectService) *cobra.Command`. El handler debe:

- validar únicamente argumentos y flags propios del comando;
- invocar un servicio mediante dependencias explícitas;
- escribir resultados en `cmd.OutOrStdout()` y errores informativos en
  `cmd.ErrOrStderr()`;
- devolver el error con `RunE`, no ocultarlo en `Run`;
- respetar `cmd.Context()` para cancelación.

Usa `Args: cobra.ExactArgs(1)`, `cobra.NoArgs` u otra validación de Cobra
cuando corresponda. Mantén `Use`, `Short`, `Long` y ejemplos útiles. Registra
completions dinámicos solo cuando aporten valor y sin hacer I/O innecesario en
cada pulsación.

Los comandos deben ser delgados. Por ejemplo, `path` pide al servicio la ruta
resuelta y la imprime; no conoce cómo se calcula la carpeta configurada. `open`
solicita al servicio una ruta validada y delega la ejecución del editor a una
abstracción de procesos.

### Servicios e infraestructura

Define la política en servicios pequeños y deja los efectos externos en
componentes concretos:

- `internal/config`: resuelve directorio raíz y editor, aplicando defaults
  mediante configuración XDG (`$XDG_CONFIG_HOME/project-cli/config.json`, con
  fallback `~/.config/project-cli/config.json`), y devuelve errores de configuración.
  También centraliza carga, validación y guardado atómico del JSON.
- `internal/services`: implementa los casos de uso de la aplicación, como
  listar, crear, resolver rutas y abrir proyectos. Orquesta configuración,
  validación, filesystem y procesos mediante dependencias explícitas. Se
  mantiene un `ProjectService` cohesivo; separar cada método en un servicio
  distinto requiere una razón comprobable.
- `internal/process`: ejecuta el editor y otros procesos mediante
  `exec.CommandContext`; conecta stdin, stdout y stderr de forma explícita.

El comando `init` debe ser el único lugar que coordina la interacción con el
usuario para configurar carpeta y editor. La configuración persistida debe
validarse antes de construir servicios; si falta o es inválida, el arranque
debe ejecutar ese flujo automáticamente sin crear un bucle de inicialización.

Los comandos consumen interfaces de servicio orientadas a casos de uso, como
`List`, `Create`, `Path` y `Open`; no resuelven rutas, leen el filesystem ni ejecutan
procesos.
Una interfaz es útil en el borde que se prueba o sustituye, por ejemplo un
servicio de aplicación o un runner de procesos. No introduzcas una interfaz
para cada struct. Las funciones deben recibir dependencias y valores; evita
estado global mutable.

## Reglas de comportamiento para `project-cli`

- Acepta solo nombres de proyecto seguros y no permitas que un nombre escape
  del directorio raíz con `..`, separadores o rutas absolutas.
- Resuelve `$HOME` mediante `os.UserHomeDir` y permite una configuración
  explícita si el producto la define; devuelve errores si no puede determinarse.
- Usa `os.MkdirAll` con permisos deliberados al crear directorios y no borres
  contenido existente sin una opción explícita.
- `create` debe rechazar proyectos existentes, crear un directorio vacío y
  abrir el editor solo después de que la creación haya tenido éxito.
- Comprueba entradas y errores de filesystem antes de abrir o ejecutar.
- Para `open`, usa el editor configurado y pasa la ruta como argumento sin
  construir un comando de shell. Nunca uses `sh -c` para interpolar entradas
  del usuario.
- Diferencia errores de usuario, configuración, filesystem y proceso con
  `fmt.Errorf("...: %w", err)` para conservar `errors.Is`/`errors.As`.
- La salida normal debe ser estable y apta para composición con otros
  comandos; no mezcles mensajes de diagnóstico con stdout.
- No ocultes errores con defaults silenciosos ni con `recover`.

## Convenciones idiomáticas de Go

- Ejecuta `gofmt` sobre todos los archivos Go modificados; el formato no es una
  decisión manual. Usa `go vet` y las herramientas oficiales del módulo antes
  de finalizar.
- Usa paquetes en minúsculas, nombres cortos y precisos, y evita abreviaturas
  no idiomáticas. Los identificadores exportados deben tener comentarios
  GoDoc que comiencen con su nombre.
- Devuelve errores como último resultado, compruébalos inmediatamente y
  envuélvelos con `%w` cuando agregues contexto. No ignores errores ni uses
  mensajes ambiguos como `error occurred`.
- Usa `context.Context` como primer parámetro en operaciones bloqueables; no lo
  guardes en structs ni uses `context.Background()` dentro de servicios.
- Prefiere la biblioteca estándar (`os`, `path/filepath`, `os/exec`, `errors`,
  `fmt`, `io`, `slices`, `maps`) antes de añadir dependencias.
- Usa tipos concretos por defecto. Define interfaces en el paquete consumidor,
  con el menor número de métodos necesario para sustituir una dependencia o
  probar un borde.
- Evita getters triviales, constructores obligatorios sin invariantes,
  abstracciones ceremoniales, reflexión y variables globales mutables.
- Evita `init()` salvo registro puramente declarativo y justificado. No uses
  `panic` para errores esperables de entrada, configuración o I/O.
- Usa `io.Reader`/`io.Writer` para entradas y salidas inyectables, y conecta
  explícitamente stdin, stdout y stderr.
- Usa `filepath` para rutas del sistema, `os.Open`/`os.ReadFile` con cierres
  comprobados y permisos explícitos para archivos creados.
- Usa `defer` para liberar recursos después de comprobar que fueron adquiridos;
  si `Close` puede fallar, propaga ese error cuando sea relevante.
- Usa `t.TempDir`, `t.Setenv`, `t.Cleanup`, `t.Run` y tests de tabla. Los tests
  no deben depender del home, PATH, editor instalado, red o reloj reales.
- Mantén una sola responsabilidad por paquete, evita ciclos de importación y
  no exportes símbolos que solo necesita otro paquete interno.

## Go 1.27

El módulo declara Go `1.27.1`. Comprueba la versión instalada con `go version`
y usa como fuente de verdad las [notas oficiales de Go
1.27](https://go.dev/doc/go1.27) y el [anuncio de
Go 1.27](https://go.dev/blog/go1.27). Go 1.27 añade, entre otras cosas:

- métodos genéricos;
- inferencia de tipos de funciones en más contextos;
- inicialización directa de campos promovidos en literales de structs;
- modernizadores adicionales en `go fix`;
- mejoras de `go doc` para consultar versiones;
- consolidación de bloques `require` por `go mod tidy`.

Estas novedades no justifican complejidad en una CLI pequeña. Usa APIs nuevas
solo cuando estén disponibles en la versión declarada y mejoren claramente la
claridad. No dependas de funciones experimentales o de comportamiento
descrito en blogs no oficiales sin verificarlo en `go.dev` y en la
documentación de `pkg.go.dev`. Ejecuta `go mod tidy` después de cambiar
dependencias y conserva `go.mod`/`go.sum` reproducibles.

## Pruebas

Prueba comportamiento observable y mantén las pruebas deterministas:

- comandos: argumentos, flags, salida, errores y códigos de ejecución usando
  `cmd.SetOut`, `cmd.SetErr` y un contexto cancelable;
- configuración: defaults, variables de entorno y errores;
- proyectos: nombres válidos e inválidos, rutas, listado y creación en
  `t.TempDir`;
- procesos: un runner falso que verifique editor, argumentos y contexto, sin
  lanzar editores reales.

Evita tests que dependan de `$HOME`, del editor instalado, del orden accidental
del filesystem o de la red. Usa tests de tabla para reglas de validación y
tests de integración pequeños para el cableado de `main` cuando sea necesario.

## Flujo de implementación

1. Lee el README y confirma el contrato de usuario antes de crear paquetes.
2. Modela la operación como una función o servicio que pueda probarse sin
   Cobra.
3. Implementa la infraestructura concreta con errores envueltos y contexto.
4. Añade el comando delgado que conecta flags, servicio y salida.
5. Registra el comando desde `internal/app`, manteniendo la composición visible.
6. Añade o actualiza pruebas cercanas al comportamiento modificado.
7. Ejecuta `gofmt`, `go vet ./...`, los tests focalizados y luego
   `go test ./...`; usa `go mod tidy` solo cuando cambien dependencias.
8. Verifica manualmente `--help`, errores de argumentos y una ejecución feliz
   sin modificar datos del usuario.

## Criterios de revisión

Antes de considerar terminada una modificación, verifica que:

- el comando nuevo está registrado y aparece en `--help`;
- ningún handler contiene lógica de negocio, filesystem o procesos;
- los errores llegan al usuario y producen un fallo no exitoso;
- stdout conserva solo la salida del contrato;
- entradas de usuario no se convierten en comandos de shell;
- la operación tiene pruebas sin depender del entorno local;
- no se agregó una dependencia cuando la biblioteca estándar era suficiente;
- el código está formateado con `gofmt` y pasa `go vet`;
- no existen ciclos de importación, estado global mutable ni `init()` con
  efectos externos;
- la solución sigue siendo sencilla de eliminar o reemplazar cuando
  `project-cli` evolucione.
