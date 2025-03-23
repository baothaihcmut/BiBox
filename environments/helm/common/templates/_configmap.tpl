{{- define "common.configmap" -}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}-{{ .Values.configMap.name }}-config-map
  labels:
    {{- include "common.labels" . | nindent 4 }}
data:
    {{- toYaml .Values.configMap | nindent 4 }}
{{- end }}
