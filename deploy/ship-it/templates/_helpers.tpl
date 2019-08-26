{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "ship-it.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "ship-it.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "ship-it.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Provide standard resource labels.
*/}}
{{ define "ship-it.metadataLabels" }}
  app: {{ template "ship-it.name" . }}
  chart: {{ template "ship-it.chart" . }}
  instance: {{ .Release.Name }}
  managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Image value hack for skaffold.
*/}}
{{- define "ship-it.api.image" -}}
{{- if .Values.api.image.tag -}}
{{- printf "%s:%s" .Values.api.image.repository .Values.api.image.tag -}}
{{- else -}}
{{- .Values.api.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "ship-it.syncd.image" -}}
{{- if .Values.syncd.image.tag -}}
{{- printf "%s:%s" .Values.syncd.image.repository .Values.syncd.image.tag -}}
{{- else -}}
{{- .Values.syncd.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "ship-it.operator.image" -}}
{{- if .Values.operator.image.tag -}}
{{- printf "%s:%s" .Values.operator.image.repository .Values.operator.image.tag -}}
{{- else -}}
{{- .Values.operator.image.repository -}}
{{- end -}}
{{- end -}}