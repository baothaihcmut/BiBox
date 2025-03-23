{{- define "common.ingressroute" -}}
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}-ingressroute
  labels:
    {{- include "common.labels" . | nindent 4 }}
spec:
  entryPoints:
    - websecure
  routes:
    {{- range .Values.ingressRoutes }}
    - match: Host(`{{ .host }}`)
      kind: Rule
      services:
        - name: {{ $.Release.Name }}-{{ $.Chart.Name }}-{{ .service }}-svc
          port: {{ .port }}
          namespace: "storage-app"
    {{- end }}
  tls:
    certResolver: letsencrypt
    domains:
      {{- range .Values.ingressRoutes }}
      {{ if .tls -}}
      - main: "{{ .host }}"
      {{- end }}
      {{- end }}
{{- end}}