#for common label
{{- define "common.labels" -}}
app: {{ .Chart.Name }}
release: {{ .Release.Name }}
env: {{ .Values.global.environment | default "dev" }}
{{- end }}

#for config from env 
{{- define "common.configEnv" -}}
{{- if (and .Values.configFrom .Values.configFrom.envConfigMap) }}
envFrom:
  {{- range .Values.configFrom.envConfigMap }}
  - configMapRef:
      name: {{ . }}
  {{- end }}
{{- end }}
{{- end }}


# for config from yaml file
{{- define "common.configFile" -}}
{{- if (and .Values.configFrom .Values.configFrom.fileConfigMap ) }}
{{- range .Values.configFrom.fileConfigMap }}
- name: {{ .volumeName }}
  mountPath: {{ .mountPath }}
  subPath: {{ .subPath }}
{{- end }}
{{- else }}
[]
{{- end }}
{{- end }}

#for config volume with file
{{- define "common.configMapVolume" -}}
{{- if (and .Values.configFrom .Values.configFrom.fileConfigMap) }}
{{- range .Values.configFrom.fileConfigMap}}
- name: {{.volumeName }}
  configMap:
    name: {{ $.Release.Name }}-{{ $.Chart.Name }}-{{ .configMap }}-config-map
{{- end }}
{{- else}}
[]
{{- end }}
{{- end }}

#for deployment init container
{{- define "common.initContainers" -}}
{{- if gt (len .Values.initContainers) 0}}
initContainers:
{{- .Values.initContainers | toYaml | nindent 2 }}
{{- end }}
{{- end }}   

# for route