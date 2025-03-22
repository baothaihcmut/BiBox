{{- define "common.service" -}}
apiVersion: v1
kind: Service
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}-svc
  namespace: {{ .Values.namespace }}
  labels:
    {{- include "common.labels" . | nindent 4 }}
spec:
  selector:
    app: {{ .Chart.Name }}
  ports:
    {{- range .Values.service.ports }}
    - name: {{ .name }}
      protocol: {{ .protocol }}
      port: {{ .port }}
      targetPort: {{ .targetPort }}
    {{- end }}
  type: {{ .Values.service.type | default "ClusterIP"}}
{{- end }}

