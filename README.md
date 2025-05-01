# TFG - Ciclo completo de CI/CD con Dagger utilizando Kubernetes

> [!warning]
> Este proyecto aún está en progreso. En este README se indica todo lo necesario para probar las funcionalidades que están implementadas.

## Estructura del monorepo

Esta es la estructura general del monorepo. Se muestran únicamente los elementos más destacables.

```bash
├── packages/
│   ├── backend/
│   └── frontend/
├── mongo-init/
├── example.env
├── justfile
├── lerna.json
├── package.json
├── yarn.lock
├── docker-compose.yaml
└── Dockerfile
```

## Funcionalidades

- [x] Configuración inicial del monorepo utilizando Lerna.
- [x] Uso de Docker Compose para la orquestación y comunicación de los servicios.
- [/] Aplicación *dummy* de gestión de un zoo:
    - [x] Configuración de la base de datos (MongoDB).
    - [x] API REST para comunicación con la base de datos (Node.js + TypeScript).
    - [/] Frontend para consumir la API (Vue.js + TypeScript).
- [ ] Uso de Kubernetes para la gestión de los servicios.
- [ ] Creación de un módulo de Dagger para la realización de un ciclo completo de CI/CD.

## Requisitos

### Software

A continuación se indica el software junto con las versiones utilizadas para el desarrollo del proyecto.


| **Software** | **Version** | **Docs** |
|---|---|---|
| Git | 2.48.1 | https://git-scm.com/ |
| Just | v1.39.0 | https://github.com/casey/just |
| Docker | v27.5.1 (con `compose` habilitado) | https://www.docker.com/ |
| Node (opcional) | 23.7.0 | https://nodejs.org/en |
| Yarn (opcional) | 1.22.22 | https://yarnpkg.com/ |
| npm (opcional) | 10.9.2 | https://www.npmjs.com/ |
| Lerna (opcional) | v8.1.9 | https://lerna.js.org/ |


### Variables de entorno

Es necesario configurar el archivo `.env`. Para ello, se proporciona un `example.env` de ejemplo, el cual hay que renombrar a `.env`.

```bash
mv example.env .env
```

```env
MONGO_DATABASE=<database_name>  # required
MONGO_PORT=<container_port>     # optional (default: 27017)
MONGO_PORT_HOST=<host_port>     # optional (default: 27017)
MONGO_ROOT=<root_name>          # required
MONGO_ROOT_PASS=<root_password> # required
CR_PAT=                         # optional (para realizar acciones con el registry remoto)
```

## Cómo probarlo

> [!note]
> Comprueba los [requisitos](#requisitos) si no lo has hecho ya.

1. Clona el repositorio y accede a él.

```bash
git clone https://github.com/vieites-tfg/zoo ~/zoo
cd ~/zoo
```

2. Instala todos los paquetes necesarios. **Este paso se debe realizar antes de cualquier otra opción**.

```bash
just init
```

3. Inicia los contenedores en modo desarrollo.

```bash
just dev
```

> [!note]
> Con `just` puedes ejecutar los pasos anteriores de manera concatenada mediante el comando `just init dev`.

4. Accede a la página web en [http://localhost:5173](http://localhost:5173)

### Más posibilidades

> [!important]
> Siempre es necesario haber realizado un `just init` previamente a la ejecución de cualquier otro comando.

- Ejecuta el linter:

```bash
just lint
```

- Comprueba que se pasan los tests, tanto del backend como del frontend. **Es necesario haber lanzado la aplicación con `just dev`** para que los tests funcionen correctamente:

```bash
just test
```

- O cualquiera de los paquetes por separado:

```bash
just test_backend # just tb
# o
just test_frontend # just tf
```

- Construye las imágenes de los paquetes (frontend o backend). La versión de la imagen que se construya se tomará como "latest":

```bash
just image_build <package> # just ib <package>
```

> [!note]
> Tanto en el anterior comando como en los que se muestan a continuación, los posibles valores para `<package>` son:
> - `backend`: Únicamente el paquete del backend.
> - `frontend`: Únicamente el paquete del frontend.
> - `all`: Tanto el backend como el frontend.
>
> **No es posible** ejecutar un comando como el siguiente ejemplo: `just ib backend frontend`. Es necesario hacer `just ib all`.

- Sube las imágenes al registry de GitHub:

```bash
just image_push <package> # just ip <package>
```

> [!note]
> Para realizar las operaciones de subir elementos al repositorio remoto, es necesario tener un token válido `CR_PAT` del archivo `.env`. En el caso de no tenerlo, es obligatorio tenerla creada con un valor vacío en (`CR_PAT=`) el archivo `.env`.

- O las dos acciones anteriores al mismo tiempo:

```bash
just image_build_push <package> # just ibp <package>
```

- Se pueden subir al *registry* remoto los paquetes npm con el siguiente comando:

```bash
just pkg_remote <package> # just pr <package>
```

- En el caso de querer almacenar los paquetes npm y no tener un token de autenticación, estos se pueden guardar de manera local. Se creará un archivo comprimido para cada uno de los paquetes que se quiera guardar y se almacenará en la raíz del repositorio, en el directorio `local_packages/`. Solo hay que ejecutar el siguiente comando:

```bash
just pkg_local <package> # just pl <package>
```

- En caso de querer realizar tanto el almacenamiento remoto como local, se ejecutaría lo siguiente:

```bash
just pkg_remote_local <package> # just prl <package>
```

> [!note]
> Los comandos anteriores para subir tanto las imágenes de Docker como los paquetes npm, utilizan los scripts `image.sh` y `push_package.sh`, respectivamente.

### Para probar la API

La API está disponible en [http://localhost:3000](http://localhost:3000). Tiene definidos los siguientes *endpoints*.

| **Acción** | **endpoint** | **Funcionalidad** |
|---|---|---|
| GET | `/animals` | Obtener todos los animales |
| GET | `/animals/{id}` | Obtener un animal mediante su ID |
| POST | `/animals` | Añadir un nuevo animal |
| PUT | `/animals/{id}` | Actualizar el animal con cierto ID |

Aquí se muestran peticiones de ejemplo que puedes probar utilizando, por ejemplo, [Postman](https://www.postman.com/).

```js
// GET http://localhost:3000/animals
```

```js
// GET http://localhost:3000/animals/<UN_ID_DEL_GET_ANTERIOR>
```

```js
// POST http://localhost:3000/animals/
{
    name: "Marcus",
    species: "Tiger",
    birthday: "2010-05-16",
    genre: "male",
    diet: "Carnivore",
    condition: "Healthy"
}
```
```js
// PUT http://localhost:3000/animals/<ID_GENERADO_DEL_POST_ANTERIOR>
{
    condition: "Injured",
    notes: "Recovering from minor foot injury."
}
```
