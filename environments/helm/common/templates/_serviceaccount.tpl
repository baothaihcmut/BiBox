{{- define "common.serviceaccount" -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}-{{ .Values.serviceAccount.name }}-sa
  namespace: {{ .Release.Namespace }}
{{- end }}