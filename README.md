# TFG - Ciclo completo de CI/CD con Dagger utilizando Kubernetes

> [!warning]
> Este proyecto aún está en progreso. En este README se indica todo lo necesario para probar las funcionalidades que están implementadas.

## Estructura del monorepo

Esta es la estructura general del monorepo. Se muestran únicamente los elementos más destacables.

```bash
├── packages/
│   ├── backend/
│   │   └── Dockerfile
│   └── frontend/
│       └── Dockerfile
├── mongo-init/
├── docker-compose.yaml
├── example.env
├── lerna.json
├── package.json
└── yarn.lock
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
| Node | 23.7.0 | https://nodejs.org/en |
| Yarn | 1.22.22 | https://yarnpkg.com/ |
| npm | 10.9.2 | https://www.npmjs.com/ |
| Docker | v27.5.1 (con `compose` habilitado) | https://www.docker.com/ |
| Lerna (opcional) | v8.1.9 | https://lerna.js.org/ |
| Just  (opcional) | v1.39.0 | https://github.com/casey/just |


### Variables de entorno

Es necesario configurar el archivo `.env`. Para ello, se proporciona un `example.env` de ejemplo, el cual hay que renombrar a `.env`.

```env
MONGO_DATABASE=<database_name>  # required
MONGO_PORT=<container_port>     # optional (default: 27017)
MONGO_PORT_HOST=<host_port>     # optional (default: 27017)
MONGO_ROOT=<root_name>          # required
MONGO_ROOT_PASS=<root_password> # required
```

## Cómo probarlo

> [!note]
> Comprueba los [requisitos](#requisitos) si no lo has hecho ya.

1. Clona el repositorio y accede a él.

```bash
git clone https://github.com/vieites-tfg/zoo ~/zoo
cd ~/zoo
```

2. Instala todos los paquetes necesarios.

```bash
make init
# o usando "just"
just init
```

3. Inicia los contenedores en modo desarrollo.

```bash
make dev
# o usando "just"
just dev
```

> [!note]
> Con `just` puedes ejecutar los pasos anteriores de manera concatenada mediante el comando `just init dev`.

4. Accede a la página web en [http://localhost:5173](http://localhost:5173)

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
