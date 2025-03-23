{{- define "common.role" -}}
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}-{{ .Values.serviceAccount.name }}-role
  namespace: {{ .Release.Namespace }}
rules:
  {{ .Values.serviceAccount.rules | toYaml | nindent 2 }}
{{- end }}