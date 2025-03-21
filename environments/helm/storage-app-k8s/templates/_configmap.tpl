{{- define "common.configmap" -}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}-config
  labels:
    {{- include "common.labels" . | nindent 4 }}
data:
  config.yaml: |-
    {{- .Values.config | nindent 4 }}
{{- end }}
