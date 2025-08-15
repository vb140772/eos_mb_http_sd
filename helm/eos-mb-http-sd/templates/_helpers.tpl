{{/*
Expand the name of the chart.
*/}}
{{- define "eos-mb-http-sd.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "eos-mb-http-sd.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "eos-mb-http-sd.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "eos-mb-http-sd.labels" -}}
helm.sh/chart: {{ include "eos-mb-http-sd.chart" . }}
{{ include "eos-mb-http-sd.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "eos-mb-http-sd.selectorLabels" -}}
app.kubernetes.io/name: {{ include "eos-mb-http-sd.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "eos-mb-http-sd.serviceAccountName" -}}
{{- if .Values.eosMbHttpSd.serviceAccount.create }}
{{- default (include "eos-mb-http-sd.fullname" .) .Values.eosMbHttpSd.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.eosMbHttpSd.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the config map to use
*/}}
{{- define "eos-mb-http-sd.configMapName" -}}
{{- default (printf "%s-config" (include "eos-mb-http-sd.fullname" .)) .Values.eosMbHttpSd.config.name }}
{{- end }}
