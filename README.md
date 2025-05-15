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
| Dagger | latest | https://dagger.io |
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

Primero hay que clonar el repositorio y acceder a él.

```bash
git clone https://github.com/vieites-tfg/zoo ~/zoo
cd ~/zoo
```

### Con Dagger

Dagger es un software construido por los creadores de Docker. Permite ejecutar cualquier tipo de workflow de manera local utilizando contenedores, todo esto de manera programática, con cada vez más SDK para diferentes lenguajes. Es la herramienta más importante en este proyecto, ya que vamos a implementar un módulo de Dagger que permita realizar un ciclo completo de CI/CD, dando la posibilidad de levantar todo el entorno con Kubernetes, no solo con contenedores de Docker.

Para probar la aplicación dummy con Dagger es necesario primeramente tener [instalado Dagger](https://docs.dagger.io/install) en su última versión.

Una vez instalado, debemos movernos al directorio correspondiente al módulo de Dagger para esta aplicación:

```bash
cd dagger
```

Desde aquí podemos ejecutar cualquir comando de dagger.

1. Muestra las funciones disponibles:

```bash
dagger functions
```

Veremos varias opciones, de las cuales nos interesan `backend` y `frontend`. Estas funciones incluyen lar acciones disponibles para cada uno de los paquetes, que en este caso son las mismas. Podemos verlas con estos comandos:

```bash
dagger call backend --help
dagger call frontend --help
```

2. Levantar la aplicación:

Para ello debemos lanzar el backend por un lado y el frontend por otro. Primero ejecutamos el comando que levantará el backend de la aplicación.

> [!note]
> Siempre es necesario indicar el archivo .env, ya que es un requisito para lanzar la aplicación tener las variables mencionadas anteriormente definidas.

Este comando habilita la API en la URL `localhost:3000/animals`. Para saber más échale un vistazo al [apartado sobre la API](#prueba-la-api)

```bash
# Terminal 1
dagger --sec-env file://../.env call backend service --ports 3000:3000
```

En otra terminal vamos a levantar el frontend con este comando:

```bash
# Terminal 2
dagger --sec-env file://../.env call frontend service --ports 8080:80
```

> [!important]
> El frontend debe levantarse siempre en el puerto 8080 para poder realizar lost test end-to-end.

El comando anterior nos permite acceder a la aplicación en la URL `localhost:8080`.

De esta manera **ya tenemos levantada completamente nuestra aplicación**:
- Una base de datos MongoDB.
- Una API conectada a la base de datos.
- Un frontend que consume dicha API.

3. Pasar test y linter.

Podemos realizar esto mediante los siguientes comandos:

```bash
# Tests
dagger --sec-env file://../.env call backend test
dagger --sec-env file://../.env call frontend test --front tcp://localhost:8080
# Linter
dagger --sec-env file://../.env call backend lint
dagger --sec-env file://../.env call frontend lint
```

> [!important]
> Como se puede observar, para pasar los tests del frontend es necesario tener toda la aplicación levantada, tanto frontend como backend, ya que se tratan de tests end-to-end. Esto es porque hay que pasar como parámetro el servicio del frontend, el cual debe estar **obligatoriamente** disponible en el puerto 8080.

4. Subir imágenes de Docker y paquetes npm al registro de GitHub.

Para subir las imágenes de Docker, es necesario estar logueado en el registro, tal y como se indica en la [documentación de GitHub](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry).

Aquí, para el frontend, también **es necesario tener levantada toda la aplicación**, ya que se pasa el linter y se ejecutan los tests previamente a publicar la imagen de Docker o el paquete npm, por si hubiera algún error.

```bash
# Imágenes de Docker
dagger --sec-env file://../.env call backend publish-image
dagger --sec-env file://../.env call frontend publish-image --front tcp://localhost:8080
# Imágenes de Docker
dagger --sec-env file://../.env call backend publish-pkg
dagger --sec-env file://../.env call frontend publish-pkg --front tcp://localhost:8080
```

### Con just

1. Instala todos los paquetes necesarios. **Este paso se debe realizar antes de cualquier otra opción**.

```bash
just init
```

2. Inicia los contenedores en modo desarrollo.

```bash
just dev
```

> [!note]
> Con `just` puedes ejecutar los pasos anteriores de manera concatenada mediante el comando `just init dev`.

3. Accede a la página web en [http://localhost:8080](http://localhost:8080)

#### Más posibilidades

> [!important]
> Siempre es necesario haber realizado un `just init` previamente a la ejecución de cualquier otro comando, al menos una vez.

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
> Sigue las [instrucciones de GitHub](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry) para crear un PAT propio y poder así subir las imágenes de Docker al registro remoto. Este se usará también para subir los paquetes npm, como se indica más abajo.

- O las dos acciones anteriores al mismo tiempo:

```bash
just image_build_push <package> # just ibp <package>
```

- Se pueden subir al *registry* remoto los paquetes npm con el siguiente comando:

```bash
just pkg_remote <package> # just pr <package>
```

> [!note]
> Para realizar las operaciones de subir elementos al repositorio remoto, es necesario tener un token válido `CR_PAT` del archivo `.env`. En el caso de no tenerlo, es obligatorio tenerla creada con un valor vacío en (`CR_PAT=`) el archivo `.env`.

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

### Prueba la API

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
