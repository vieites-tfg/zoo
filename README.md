# TFG - Ciclo completo de CI/CD con Dagger utilizando Kubernetes

## Repositorios

Este trabajo se compone de tres repositorios bien diferenciados y con su protia funcionalidad. Todos ellos se pueden encontrar en esta organización, `vieites-tfg`.

- `zoo`: Este repositorio. Es el principal del trabajo. En él se encuentra todo el código fuente, desde la aplicación *dummy* hasta la implementación de los módulos de Dagger. Es el único necesario para probar las diferentes funcionalidades.
- `helm-repository`: Alberga las Chart de Helm que describen el despliegue de la aplicación *dummy* en Kubernetes.
- `state`: En él se almacenan los valores de las Charts de Helm correspondientes a cada uno de los entornos en los que se puede desplegar la aplicación. Funciona como única fuente de verdad. También es el lugar de donde ArgoCD obtiene los recursos de Kubernetes, con el fin de ser desplegados.

Cada uno de los repositorios tiene su propio `README`. Consultar los demás para más información.

# Zoo

## Estructura general del repositorio

```bash
argo
├── argo_dev.yaml
├── argo_pre.yaml
├── argo_pro.yaml
└── values.yaml
cluster
├── kind_dev.yaml
├── kind_local.yaml
├── kind_pre.yaml
└── kind_pro.yaml
dagger
├── cd
└── ci
docker-compose.yml
Dockerfile
example.env
mongo-init
└── init-zoo.js
packages
├── backend
└── frontend
scripts
├── create_envs.sh
├── image.sh
└── push_package.sh
sops
├── .sops.yaml
└── age.agekey
```

Este repositorio funciona como un monorepo, lo cual indica que todo el código fuente se encuentra en este único lugar.

### Aplicación *dummy*

Consta de un *frontend* y un *backend*, ambos creados con Typescript, y utilizando Vue para el frontend. Se gestiona utilizando una herramienta de gestión de monorepos llamada Lerna. En la raíz del repositorio se pueden encontrar archivos de configuración de Lerna y Node.js. El código correspondiente al *frontend* y al *backend* se incluyen en el directorio `packages`. El *backend* se conecta a una base de datos de animales de un zoo, y proporciona una API que el *frontend* consume, con el fin de poder realizar acciones CRUD sobre la base de datos.

### Dagger

En el directorio `dagger` se encuentran las implementaciones de los módulos correspondientes a los flujos de CI y CD. Estos se pueden ejecutar localmente teniendo Dagger y Docker instalados. Se utilizan en el workflow de GitHub encargado de realizar todo el flujo de testeo, publicación y despliegue.

### Otros

También se pueden encontrar varios ejecutables en el directorio de `scripts`, de los cuales `create_envs.sh` es el más interesante. Este permite levantar los clusters en local para probar el despliegue de la aplicación en los tres entornos posibles: `dev`, `pre` y `pro`.

En todos los clusters se instala ArgoCD, aplicación creada específicamente para Kubernetes y que utiliza el método *pull*, siguiendo la filisofía GitOps, leyendo los recursos a desplegar del repositorio de estado (`state`) mencionado anteriormente. Este tiene una rama `deploy`, en la que se suben los archivos necesarios para que Argo lea y despliegue la aplicación.

Se configuran los clusters y las propias instancias de Argo de manera diferente para cada uno de los clusters. Estas configuraciones se pueden encontrar en los directorios `cluster` y `argo`.

Además, este trabajo tiene una memoria asociada, en la cual se pueden encontrar unas claves privada y pública, junto con otros *tokens*, necesarias para poder probar la implementación. Estas claves y *tokens* tienen varios propósitos:

- Claves privada y pública: Se utilizan para encriptar y desencriptar los secretos (`Secrets`) de Kubernetes.
- `tokens`: Un `CR_PAT`, que permite acceder al Container Registry de GitHub, donde se encuentran las imágenes de Docker de la aplicación; y un `STATE_REPO`, necesario para poder actualizar el repositorio de estado y actualizar los recursos que se van a desplegar.

## Prueba mínima

Aquí se describen los requisitos de software y pasos a seguir para probar de la manera más simple la aplicación. Con esta prueba se comprobará que se pueden obtener los recursos de Kubernetes y levantar la aplicación en los diferentes entornos, pudiendo visualizar todos estos recursos en ArgoCD.

Los elementos que van a influir en esta prueba serán: los clusters junto con sus configuraciones y las de ArgoCD, y el repositorio de estado, en el cual ya debería haber recursos preparados para desplegar.

### Requisitos software

A continuación se indica el software junto con las versiones utilizadas para el desarrollo del proyecto.

> [!note]
> **No** se ha probado en un sistema operativo Windows, por lo que no se asegura su funcionamiento en este.
> Sí se ha probado en MacOS y distrybuciones Linux.

| **Software** | **Version** | **Docs** |
|---|---|---|
| Git | 2.48.1 | https://git-scm.com/book/en/v2/Getting-Started-Installing-Git |
| Just | v1.39.0 | https://github.com/casey/just?tab=readme-ov-file#installation |
| Docker | v27.5.1 | https://docs.docker.com/desktop/ |
| Kubectl | v1.33 | https://kubernetes.io/docs/tasks/tools/#kubectl |
| Kind | 0.27.0 | https://kubernetes.io/docs/tasks/tools/#kind |
| Helm | v3.17.3 | https://helm.sh/docs/intro/install/ |

### Cómo probar

1. Clonar este repositorio.

Se clona el repositorio y se accede al directorio.

```bash
git clone https://github.com/vieites-tfg/zoo ~/zoo
cd ~/zoo
```

2. Clave privada

> [!warning]
> Reservar únicamente para la entrega final de la memoria.

```bash
mkdir -p ./sops # crear el directorio en el caso de que no exista
echo "AGE-SECRET-KEY-1CTS4S4QNNZ9N9YXXM288LSE9VKPJ220E57ZHC4558WMZ8LG2QWKQFFER8C" > ./sops/age.agekey
```

3. Ejecutar el script de creación de los entornos.

```bash
./scripts/create_envs.sh
```

El script anterior:
- Crea tres clusters (`dev`, `pre` y `pro`), con tres contextos diferentes (`kind-{{cluster}}`), con su propia configuración.
- Introduce la clave privada, creada previamente, en cada uno de los clusters, para permitir a ArgoCD desencriptar los secretos.
- Instala ArgoCD en cada uno de los clusters, con sus respectivas configuraciones, obteniendo cada uno los recursos de despliegue del entorno que le toca.

4. Acceso a los clusters.

A medida que se van creando los clusters, las contraseñas del usuario `admin` de Argo se van mostrando. También se muestran todas al finalizar la ejecución del script.

Para poder acceder a cada uno de los clusters, lo primero que hay que hacer es mapear un puerto local libre al puerto 443 del servidor de Argo. Esto se consigue de la siguiente manera:

```bash
kubectl port-forward svc/argocd-server -n argocd --context kind-{{cluster}} 8086:443
```

En el anterior comando, hay que cambiar `{{cluster}}` por aquel al que se quiera acceder. Se podrá acceder a Argo a través del navegador en `localhost:8086`

Se pide usuario y contraseña para entrar, que son `admin` y la contraseña de dicho cluster, mostrada en la salida del script que se ha ejecutado antes.

5. Acceder a la aplicación. (opcional)

Lo primero que hay que hacer es configurar los hosts del ordenador para que se resuelvan las rutas como `localhost`. Para ello se puede ejecutar el siguiente comando:

```bash
just check_hosts dev pre pro
```

El comando anterior modifica el archivo `/etc/hosts`, incluyendo las líneas nesecarias para poder resolver las rutas de acceso a la aplicación. Es necesario tener permisos de usuario o poner la contraseña de este en el caso de que se pida.

Hay que tener en cuenta que cada entorno tiene su propio puerto, que son:
- `dev`: 8080
- `pre`: 8081
- `pro`: 8082

> [!note]
> Los anteriores puertos se podrían modificar, pero sería necesario actualizar tanto las configuraciones de los clusters en `zoo/cluster` como los valores de los puertos del Ingress en el repositorio de estado `state` para cada uno de los entornos. Esto implicaría tener que hacer un nuevo despliegue con el fin de actualizar los valores en los propios recursos de Kubernetes. No sería necesario publicar nuevas imágenes.

Ahora se puede acceder a la aplicación de gestión del zoo. a través de las siguientes rutas en el navegador:

- *frontend*: `zoo-{{entorno}}.example.com:{{puerto_entorno}}`
- *backend*: `api-zoo-{{entorno}}.example.com:{{puerto_entorno}}`

### Conclusión

Se comprueba que funciona el despliegue de la aplicación para cualquiera de los entornos. Esta configuración permite al desarrollador tener a su disposición cada una de las versiones y comprobar que el despliegue se realiza correctamente. El uso de ArgoCD y la capacidad de obtener los recursos de una única fuente de verdad hace de este flujo de despliegue algo esencial en cualquier equipo de desarrollo.

Es necesario mencionar que, en entornos de producción reales, los clusters se encontrarían en la nube. Sin embargo, se podría mantener el cluster de desarrollo `dev`, con el fin de realizar implementaciones y probar su funcionamiento sin depender de herramientas remotas.

A continuación se indica cómo probar los módulos de Dagger.

## Prueba de Dagger

En este caso se va mostrar cómo probar tanto los módulos de Dagger como el workflow completo de CI/CD, desde la simulación de implementación de una nueva característica en la aplicación, pasando por su despliegue en todos los entornos.

### Requisitos de software

Incluyendo los de la prueba anterior.

| **Software** | **Version** | **Docs** |
|---|---|---|
| Dagger | latest | https://docs.dagger.io/install/ |
| act | 0.2.79 | https://nektosact.com/installation/index.html |

Para la simulación de una nueva *feature* se utilizar *act*, herramienta que permite ejecutar workflows de GitHub de manera local. Con ella se simulará el disparo de los eventos que hac que el workflow se ejecute, y se comportará de la misma manera que haría en remoto.

### Cómo probar

1. Creación de archivos `.env` y `.secrets.yaml`

Es necesario modificar los archivos `example.env` y `example.secrets.yaml` con los datos pertinentes.

Así quedarían los archivos:

```bash
# .env
MONGO_DATABASE=zoo
MONGO_ROOT=carer
MONGO_ROOT_PASS=carerpass
CR_PAT=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
STATE_REPO=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
```

```yaml
# .secrets.yaml
MONGO_DATABASE: zoo
MONGO_ROOT: carer
MONGO_ROOT_PASS: carerpass
CR_PAT: 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
STATE_REPO: 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
SOPS_CONFIG_FILE: |
  creation_rules:
      - path_regex: ".*\\\\.ya?ml$"
        unencrypted_regex: "^(apiVersion|metadata|kind|type)$"
        age: age15peyc7pedj8gjqwnarat6s3u87wy4j5xtf7t96vuj74m3l9xq5ys0r4sag
SOPS_PRIVATE_KEY: AGE-SECRET-KEY-1CTS4S4QNNZ9N9YXXM288LSE9VKPJ220E57ZHC4558WMZ8LG2QWKQFFER8C
```

Se puede copiar el contenido anterior y pegarlo en su respectivo archivo. Posteriormente, hay que cambiar el nombre de estos, eliminando la parte de `example`.

```bash
mv example.env .env
mv example.secrets.yaml .secrets.yaml
```

2. Prueba local del módulo de CI.

El módulo de CI de Dagger se divide en funciones para *frontend* y para *backend*, por separado. Las funciones son las mismas para ambos, pero internamente se comportan diferente.

Aquí se muestra un diagrama de la implementación de este módulo.

![ci_schema](assets/ci_schema.png)

En el diagrama se ve:
- `CI`: La estructura principal, con cinco funciones, dos de las cuales permiten acceder a las estructuras de *frontend* y *backend*, por separado.
- `Frontend` y `Backend`: Estructuras dedicadas, con implementación diferente para cada una de las funciones que proporcionan, que en este caso son las mismas para ambos.

Para poder utilizar el comando `dagger`, es necesario estar en un directorio de trabajo en el que exista un módulo de Dagger, o bien proporcionarlo con la opción `-m`. Para facilitar la explicación, se ejecutarán los comandos desde el directorio correspondiente al módulo de CI.

```bash
cd dagger/ci
```

Se pueden obtener las funciones de `CI`, junto con los campos de la estructura, con el comando:

```bash
dagger functions
```

Para conocer las funciones y campos de las demás estructuras se ejecuta:

```bash
dagger call backend --help
dagger call frontend --help
```

`CI` tiene un parámetro requerido, que se trata del archivo `.env` que se ha creado anteriormente. A continuación se muestra cómo se ejecutarían los tests *end-to-end* de la aplicación, teniendo en cuenta que hay que encontrarse en el directorio del módulo y que el archivo `.env` se ha creado correctamente en la raíz del repositorio:

```bash
dagger call --sec-env=file://../../.env endtoend
```

La primera vez que se ejecuta el comando anterior puede tardar alrededor de 10 minutos. Siempre dependiendo de la conexión a Internet que se tenga. La segunda vez, ese tiempo debería reducirse considerablemente, alrededor de un 40%, tardando así 6 minutos. Este tiempo puede reducirse más, a medida que se realizan ejecuciones del mismo, hasta alrededor de un 60% del tiempo inicial.

Otro ejemplo sería, levantar el frontend y el backend y hacer que se comuniquen de manera local.

```bash
dagger call --sec-env=file://../../.env backend service up --ports 3010:3000
dagger call --sec-env=file://../../.env frontend service up --ports 8090:80
```

Simplemente, ejecutando los comandos anteriores en terminales diferentes, los servicios serán capaces de comunicarse. Estos estarán disponibles en `localhost:{{puerto}}`. Para acceder a la API se añade la ruta `/animals`.

En los comandos anteriores, `up` es una función del tipo Service, propio de Dagger, que se devuelve en la función `service`, como se muestra en el diagrama anteriro.`--ports` es un argumento de la función `up`.

Otros ejemplos de comandos:

```bash
dagger call --sec-env=file://../../.env [backend|frontend] lint
dagger call --sec-env=file://../../.env [backend|frontend] test
dagger call --sec-env=file://../../.env [backend|frontend] publish-image --tag "{{tag}}"
dagger call --sec-env=file://../../.env [backend|frontend] publish-pkg
```

3. Prueba local del módulo de CD.

El módulo de CD tiene la siguiente implementación.

![cd_schema](assets/cd_schema.png)

La función principal de este módulo es `deploy`, que se encarga de:
- Construir todos los recursos de Kubernetes, haciendo uso de los repositorios `helm-repository` y `state`.
- Encriptar los secretos.
- Subir los cambios al repositorio de estado.

Para conseguir esto, se hace uso de [este módulo](https://daggerverse.dev/mod/github.com/prefapp/daggerverse/kind@42985961eb3d61fa98aa71d2f67922a933b5caa3), que permite crear un cluster de KinD.

> [!note]
> No se recomienda la ejecución en local de este módulo, pero se explica igual cómo se haría.

Para realizar un despliegue habría que ingresar al directorio del módulo y ejecutar el siguiente comando:

```bash
cd ~/zoo/dagger/cd

just deploy dev
# just deploy pre
# just deploy pro

# lo anterior ejecuta:
dagger call \
    --socket=/var/run/docker.sock \
    --kind-svc=tcp://localhost:3000 \
    --config-file=file://../../cluster/kind_local.yaml \
    launch \
    --sec-env=file://../../.env \
    --env={{env}} # <- cambiar {{env}} por "dev", "pre" o "pro" \
    --age-key=file://../../sops/age.agekey \
    --sops-config=file://../../sops/.sops.yaml
```

3. Prueba con `act`.

`act` es una herramienta que permite ejecutar workflows de GitHub en local, pudiendo indicar el *trigger* que dispara el workflow. 

De esta manera, se puede probar cómo sería el flujo de ejecución en el caso de que se introdujera en la rama principal una nueva característica de la aplicación.

El workflow encargado de realizar las operaciones de CI y CD es el que se encuentra en `.github/workflows/cicd.yaml`.

La definición de los objetos del *trigger* `release` se encuentran en el directorio `.github/workflows/events`.

Una vez se tenga `act` instalado, se puede iniciar la **promoción de entornos**.

Antes de nada, para comprobar que se ha realizado el cambio, se realizarán los siguientes pasos para asegurarse de que las *tags* de la imagen de `dev` es nueva:

```bash
git switch -c environments-test
```

El comando anterior creará una rama nueva con el nombre `environments-test`.

Ahora se puede cambiar el título de la página web. Este se encuentra en el archivo `packages/frontend/src/components/MainTitle.vue`, de "Zoo" a algo como "Mi zoo".

Se debe crear un *commit* para que la imagen use como *tag* los ocho primeros caracteres de este:

```bash
git add .
git commit -m 'cambio de nombre'
```

Así se utilizará el *commit* anterior para generar el nombre de la *tag* de las imágenes.

#### dev

Ahora se puede subir todo al entorno de "dev", simulando un `push` (*trigger* por defecto de `act`) a la rama principal, como si el cambio anterior se hubiera realizado en una PR y si hubieran aprobado los cambios. Esto se puede hacer, desde la raíz del repositorio, con el siguiente comando:

```bash
# se utiliza el archivo .secrets.yaml creado anteriormente
act --secret-file .secrets.yaml
```

Una vez terminada la ejecución, si se tiene levantado el cluster de `dev`, se debería poder sincronizar Argo. Para comprobar que se ha actualizado la imagen correctamente, se puede visualizar la *tag* que está en uso en el Deployment del frontend `zoo-dev-frontend`, y buscar el campo `spec.template.spec.containers.image`, en la pestaña de "*Live manifest*". Este campo debería contener la nueva imagen generada.

También se puede entrar en la propia web, en `zoo-dev.example.com:8080`, para ver el nuevo título.

#### pre

Para pasar los cambios que se han realizado al entorno de `pre`, sería neceario crear una `pre-release`. Se va a simular su creación utilizando la configuración de dicho evento, que se encuentra en `.github/events/prerelease.json`.

Lo mejor para comprobar que se actualiza correctamente, es cambiar el nombre de la `pre-release`, en el campo `name` del archivo del evento. Por ejemplo, `0.0.123-snapshot`.

Con el siguiente comando se realizará la simulación:

```bash
act release --eventpath .github/events/prerelease.json --secret-file .secrets.yaml
```

El archivo de comprobación del recurso es, en este caso, `zoo-pre-frontend`, y el mismo campo que antes.

La url sería `zoo-dev.example.com:8081`.

#### pro

Ahora, es necesario realizar algo parecido al paso anterior. Se modifica el archivo `.github/events/release.json`, y se pone el mismo `name` que anter, pero sin la coletilla `snapshot`, es decir, `0.0.123`.

Se simula el cambio al entorno de `pro`.

```bash
act release --eventpath .github/events/release.json --secret-file .secrets.yaml
```

### Conclusión

Se puede comprobar que el ciclo completo funciona correctamente, y la flexibilidad que da poder ejecutar los módulos de Dagger en local. Esto facilita el desarrollo de características nuevas para la aplicación sin miedo a que se produzca un fallo durante la ejecución del workflow de GitHub en la nube. Así, el resultado que va a dar la ejecución de dicho workflow ya se sabe de antemano, porque se ha probado de manera local, en un entorno controlado gracias a Docker. Por lo tanto, se sabe que este se va a comportar de la misma manera en cualquier otro sistema.


### Con Kubernetes

Para probarlo necesitaremos tener Kind, Helm y helmfile instalados. Se proporcionan más arriba los links correspondientes.

Ejecutaremos los siguientes comandos:

> [!important]
> Es necesario tener un CR_PAT válido almacenado en el `.env` para ser capaces de obtener las imágenes del frontend y del backend.

En el caso de ser la **primera vez que se realiza la instalación de la chart** en local, es necesario instalar el repositorio de [bitnami](https://github.com/bitnami/charts), del cual se utiliza la [chart de mongodb](https://artifacthub.io/packages/helm/bitnami/mongodb/16.5.12). Esto se puede hacer con el siguiente comando:

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami`
```

A continuación, ejecutamos los comandos necesarios para levantar el cluster de Kind con posibilidad de incluir un ingress y aceptar así conexiones desde el exterior. Entonces, levantamos la chart en el namespace que queramos.

```bash
just create_cluster # Crea el cluster de kind y los namespaces
just apply_ingress launch_chart dev # Aplica un ingress y lanza la chart
```

La chart se lanza obteniendo los valores del archivo `.env`, y se levanta el namespace indicado.

Una vez ejecutados los comandos anteriores, tendremos a nuestra disposición las urls necesarias para acceder tanto al frontend como al backend.

Podemos acceder a cada uno de ellos en:

```bash
http://zoo-[namespace].example.com:8080             # frontend
http://api-zoo-[namespace].example.com:8080/animals # backend
```

Con [k9s](https://k9scli.io/) podemos comprobar todo lo que se ha levantado.

Para desinstalar la release de helm que acabamos de instalar, podemos ejecutar:

```bash
helm uninstall zoo-[namespace] -n [namespace]
```

### Sin Kubernetes

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
