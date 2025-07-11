# Helm-repository

Este repositorio forma parte de un TFG, donde el repositorio principal es [zoo](http://github.com/vieites-tfg/zoo).

En este reposito se albergan las Charts Helm correspondientes a la aplicación que se implementa en el repositorio principal `zoo`.

A continuación se muestra un diagrama de la disposición de las Charts.

![charts diagram](assets/charts_diagram.png)

Como se puede comprobar, existe una *Chart Umbrella* `zoo`, la cual se encarga de indicar las subcharts a crear. Entre estas hay dos creadas específicamente para la aplicación, `zoo-backend` y `zoo-frontend`, perfectamente autodescritas. Por último está la Chart de MongoDB, obtenida del repositorio de [Bitnami](https://bitnami.com/stacks?stack=helm).

## Estructura del repositorio

```bash
helm
├── zoo-0.0.0.tgz
├── zoo-0.0.1.tgz
├── zoo-0.0.2.tgz
└── zoo-0.0.3.tgz
zoo
├── Chart.yaml
├── charts
│   ├── zoo-backend
│   │   ├── Chart.yaml
│   │   └── templates
│   │       ├── configMap.yaml
│   │       ├── deployment.yaml
│   │       ├── ingress.yaml
│   │       ├── secret.yaml
│   │       └── service.yaml
│   │
│   └── zoo-frontend
│       ├── Chart.yaml
│       └── templates
│           ├── configmap-js.yaml
│           ├── configmap.yaml
│           ├── deployment.yaml
│           ├── ingress.yaml
│           └── service.yaml
│   
└── templates
    ├── _helpers.tpl
    └── ghcr-secret.yaml
```


