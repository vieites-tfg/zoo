{{- define "dockerconfigjson" -}}
{
  "auths": {
    "{{ .Values.global.ghcrSecret.registry }}": {
      "username": "{{ .Values.global.ghcrSecret.username }}",
      "password": "{{ .Values.global.ghcrSecret.password }}",
      "auth": "{{ printf "%s:%s" .Values.global.ghcrSecret.username .Values.global.ghcrSecret.password | b64enc }}"
    }
  }
}
{{- end -}}
