{{- define "common.labels" -}}
app: {{ .Chart.Name }}
release: {{ .Release.Name }}
env: {{ .Values.global.environment | default "dev" }}
{{- end }}


{{- define "common.configEnv" -}}
{{- if .Values.configFrom.envConfigMap }}
envFrom:
  {{- range .Values.configFrom.envConfigMap }}
  - configMapRef:
      name: {{ . }}
  {{- end }}
{{- end }}
{{- end }}

{{- define "common.configFile" -}}
{{- if .Values.configFrom.fileConfigMap }}
{{- range .Values.configFrom.fileConfigMap }}
- name: {{ .volumeName }}
  mountPath: {{ .mountPath }}
  subPath: {{ .subPath }}
{{- end }}
{{- else }}
[]
{{- end }}
{{- end }}

{{- define "common.configMapVolume" -}}
{{- if .Values.configFrom.fileConfigMap }}
{{- range .Values.configFrom.fileConfigMap}}
- name: {{.volumeName }}
  configMap:
    name: {{ .configMap }}
{{- end }}
{{- else}}
[]
{{- end }}
{{- end }}