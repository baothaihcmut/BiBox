{{- define "common.ingress" -}}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}-{{ .Values.ingress.name }}-ingress
  labels:
    {{- include "common.labels" . | nindent 4 }}
  {{- with .Values.ingress.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  ingressClassName: {{ .Values.ingress.className | default ""}}
  rules:
    {{range .Values.ingress.rules }}
    - host: {{ .host }}
      http:
        paths:
          - path: {{ .path }}
            pathType: Prefix
            backend:
              service:
                name: {{ $.Release.Name }}-{{ $.Chart.Name }}-{{ .service }}-svc
                port: 
                  number: {{ .port }}
    {{end}}
  tls:
    {{range .Values.certs }}
    - hosts: {{ .hosts }}
      secretName: {{ .secretName | default (printf "%s-tls" .host) }}
    {{- end }}
{{- end }}