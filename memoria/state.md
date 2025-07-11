# state

Este repositorio forma parte de un TFG, donde el repositorio principal es [zoo](https://github.com/vieites-tfg/zoo).

En este repositorio se almacenan los valores necesarios para las Charts de Helm correspondientes a la aplicación, que se encuentran en [helm-repository](https://github.com/vieites-tfg/helm-repository).

## Estructura del repositorio

```bash
dev
├── global.yaml
├── mongodb.yaml
├── zoo-backend.yaml
└── zoo-frontend.yaml
pre
├── global.yaml
├── mongodb.yaml
├── zoo-backend.yaml
└── zoo-frontend.yaml
pro
├── global.yaml
├── mongodb.yaml
├── zoo-backend.yaml
└── zoo-frontend.yaml
dev.yaml
global.yaml
helmfile.yaml.gotmpl
mongodb.yaml
pre.yaml
pro.yaml
zoo-backend.yaml
zoo-frontend.yaml
```

Se ven los siguientes archivos:
- `global`: En estos se indican valores globales a toda la aplicación.
- `zoo-backend`, `zoo-frontend` y `mongodb`: Archivos de valores de cada una de las Subcharts.
- `helmfile.yaml.gotmpl`: Archivo en el que se indican los repositorios a utilizar y los valores a incluir en las Charts. Se utiliza con la herramienta [helmfile](https://helmfile.readthedocs.io/en/latest/) para generar la plantilla de todos los recursos que se van a construir.

Existen directorios para cada uno de los posibles entornos, con valores específicos de dicho entorno para cada una de las Subcharts o valores globales.

## Rama de despliegue

La rama de despliegue `deploy` es la rama en la que se publican los recursos que ArgoCD va a leer para construir la aplicación.

### Estructura

```bash
dev
├── kustomization.yaml
├── non-secrets.yaml
├── secret_generator.yaml
└── secrets.yaml
pre
├── kustomization.yaml
├── non-secrets.yaml
├── secret_generator.yaml
└── secrets.yaml
pro
├── kustomization.yaml
├── non-secrets.yaml
├── secret_generator.yaml
└── secrets.yaml
```

Se puede comprobar que hay un directorio con los recursos correspondientes de cada uno de los entornos posibles.

En ellos encontramos:

- `non-secrets`: Los recursos de Kubernetes que no son secretos.
- `secrets`: Los recursos de Kubernetes que son secretos, encriptados utilizando [SOPS](https://github.com/getsops/sops) y [age](https://github.com/FiloSottile/age).
- `kustomization`: Archivo que lee la herramienta [kustomize](https://kustomize.io/) e indica los recursos y el generador a utilizar para crear los secretos desencriptados.
- `secret_generator`: Indica la necesidad del uso de la herramienta [ksops](https://github.com/kubernetes/kops) con el fin de desencriptar los secretos.

