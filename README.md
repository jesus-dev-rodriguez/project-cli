# pjt

`pjt` es una CLI escrita en Go para gestionar y abrir proyectos locales
ubicados en `~/Proyectos`. Incluye autocompletado dinámico para facilitar el
uso desde la terminal.

La aplicación usa `cmd/project-cli` como entrypoint y mantiene su implementación
privada en `internal/`. Esta organización es una convención adecuada para un
repositorio con un ejecutable; Go no exige que todos los proyectos tengan esas
dos carpetas.

## Funcionalidades

- `pjt list`: lista los proyectos disponibles.
- `pjt open <proyecto>`: abre un proyecto con el editor configurado.
- `pjt init`: configura la carpeta de proyectos y el editor preferido.
- `pjt create <proyecto>`: crea un proyecto nuevo y lo abre con el editor
  configurado.
- `pjt path <proyecto>`: muestra la ruta absoluta de un proyecto.
- `pjt completion <shell>`: genera el autocompletado para `bash`, `zsh` o
  `fish`.

## Instalación

### Requisitos

- ⚠️ Solo compatible con Linux.
- ⚠️ Go instalado.

### Compilar e instalar

Clona el repositorio y entra en la carpeta del proyecto:

```bash
git clone <URL_DEL_REPOSITORIO>
cd hola-cli
```

Compila el binario y guárdalo en `~/.local/bin`:

```bash
mkdir -p ~/.local/bin
go build -ldflags="-s -w" -o ~/.local/bin/pjt ./cmd/project-cli
```

Asegúrate de que `~/.local/bin` esté en tu `PATH`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Para conservar este cambio, añade la línea anterior a `~/.bashrc` o
`~/.zshrc`. Comprueba la instalación con:

```bash
pjt --help
```

### Releases

Los releases se generan automáticamente al publicar un tag con formato `v*`.
El workflow compila un binario Linux `amd64` sin CGO y lo publica como asset:

```bash
git tag v0.1.0
git push origin v0.1.0
```

El workflow utiliza un runner hospedado por GitHub (`ubuntu-latest`) y configura
Go 1.27.1 automáticamente. GitHub CLI (`gh`) ya está disponible en la imagen
del runner, y el workflow requiere el permiso `contents: write`.

También puedes ejecutar el workflow manualmente sobre un tag existente, sin
crear una nueva versión:

```bash
gh workflow run release.yml -f release_tag=v1.0.0
gh run watch
```

El tag indicado debe existir en GitHub. Esta opción recompila el binario y
actualiza el asset del release con `--clobber`.

### Autocompletado

#### Bash

Genera el script:

```bash
pjt completion bash > ~/.local/share/pjt-completion.bash
```

Carga el script desde `~/.bashrc`:

```bash
source ~/.local/share/pjt-completion.bash
```

Después, recarga la configuración:

```bash
source ~/.bashrc
```

#### Zsh

Genera el script:

```bash
mkdir -p ~/.zsh/completions
pjt completion zsh > ~/.zsh/completions/_pjt
```

Añade `~/.zsh/completions` al `fpath` de tu `~/.zshrc` y recarga la
configuración:

```bash
autoload -Uz compinit && compinit
source ~/.zshrc
```

#### Fish

Genera el script directamente en la carpeta de completados de Fish:

```bash
mkdir -p ~/.config/fish/completions
pjt completion fish > ~/.config/fish/completions/pjt.fish
```

## Uso

Abre un proyecto:

```bash
pjt open mi-app
```

Crea un proyecto:

```bash
pjt create nueva-api
```

El comando crea un directorio vacío dentro de la carpeta configurada y abre
automáticamente el editor en esa carpeta. Si el proyecto ya existe, termina
con un error y conserva su contenido.

Obtén la ruta de un proyecto para cambiar de directorio:

```bash
cd "$(pjt path mi-app)"
```

`pjt path <proyecto>` valida que el proyecto exista y muestra únicamente su
ruta absoluta, por lo que puede componerse con otros comandos.

La primera vez que ejecutes `pjt`, o cuando su configuración no exista o no
sea válida, mostrará automáticamente el asistente de configuración. También
puedes ejecutarlo manualmente para cambiar estos valores:

```bash
pjt init
```

La configuración se guarda como JSON en
`$XDG_CONFIG_HOME/project-cli/config.json`. Si `XDG_CONFIG_HOME` no está
definido, se usa `~/.config/project-cli/config.json`. Por defecto, la carpeta
de proyectos propuesta es `~/Proyectos`.

El archivo contiene la carpeta de proyectos y el editor seleccionado:

```json
{
  "projects_dir": "/home/usuario/Proyectos",
  "editor": "code"
}
```

La carpeta se crea automáticamente durante la inicialización. Si el archivo
existe pero contiene JSON inválido, apunta a una carpeta inexistente o usa un
editor no disponible en el `PATH`, `pjt` vuelve a solicitar la configuración.

## Editor predeterminado

`pjt init` muestra los editores conocidos instalados (`code`, `cursor`, `zed`,
`nvim`, `vim` y `emacs`) para que selecciones uno sin escribir el comando.
Solo si no detecta ninguno permite indicar manualmente un ejecutable disponible
en el `PATH`.
