{{- define "common.deployment" -}}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}-deployment
  namespace: {{ .Values.namespace }}
  labels:
    {{- include "common.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      app: {{ .Chart.Name }}
  template:
    metadata:
      labels:
        {{- include "common.labels" . | nindent 8 }}
    spec:
      {{- include "common.initContainers" . | nindent 6}}        
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            {{- range .Values.service.ports }}
            - name: {{ .name }}
              containerPort: {{ .targetPort }}
            {{- end }}
          {{- include "common.configEnv" . | nindent 10}}
          volumeMounts:
            {{- include "common.configFile" . | nindent 12 }}
          resources:
            requests:
              cpu: {{ .Values.resources.requests.cpu }}
              memory: {{ .Values.resources.requests.memory }}
            limits:
              cpu: {{ .Values.resources.limits.cpu }}
              memory: {{ .Values.resources.limits.memory }}
      volumes:
        {{- include "common.configMapVolume" . | nindent 8}}
{{- end }}

