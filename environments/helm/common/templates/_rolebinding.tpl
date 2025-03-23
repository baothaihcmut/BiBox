{{- define "common.rolebinding" -}}
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}-{{ .Values.serviceAccount.name }}-rolebinding
  namespace: {{ .Release.Namespace }}
subjects:
  - kind: ServiceAccount
    name: {{ .Release.Name }}-{{ .Chart.Name }}-{{ .Values.serviceAccount.name }}-sa
    namespace: {{ .Release.Namespace }}
roleRef:
  kind: Role
  name: {{ .Release.Name }}-{{ .Chart.Name }}-{{ .Values.serviceAccount.name }}-role
  apiGroup: rbac.authorization.k8s.io
{{- end }}