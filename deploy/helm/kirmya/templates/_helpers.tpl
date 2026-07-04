{{/* Chart name, optionally overridden. */}}
{{- define "kirmya.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Fully-qualified app name. */}}
{{- define "kirmya.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name (include "kirmya.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/* Common labels. */}}
{{- define "kirmya.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
app.kubernetes.io/name: {{ include "kirmya.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: kirmya
{{- end -}}

{{/* Selector labels for a component. Usage: include "kirmya.selectorLabels" (dict "ctx" . "component" "backend") */}}
{{- define "kirmya.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kirmya.name" .ctx }}
app.kubernetes.io/instance: {{ .ctx.Release.Name }}
app.kubernetes.io/component: {{ .component }}
{{- end -}}

{{/* Component resource name. Usage: include "kirmya.componentName" (dict "ctx" . "component" "backend") */}}
{{- define "kirmya.componentName" -}}
{{- printf "%s-%s" (include "kirmya.fullname" .ctx) .component -}}
{{- end -}}

{{/* In-cluster service hostnames. */}}
{{- define "kirmya.postgresHost" -}}{{ include "kirmya.fullname" . }}-postgres{{- end -}}
{{- define "kirmya.redisHost" -}}{{ include "kirmya.fullname" . }}-redis{{- end -}}
{{- define "kirmya.opensearchHost" -}}{{ include "kirmya.fullname" . }}-opensearch{{- end -}}
{{- define "kirmya.natsHost" -}}{{ include "kirmya.fullname" . }}-nats{{- end -}}
{{- define "kirmya.backendHost" -}}{{ include "kirmya.fullname" . }}-backend{{- end -}}

{{/* DATABASE_URL — composed from the in-cluster Postgres, or the external override. */}}
{{- define "kirmya.databaseUrl" -}}
{{- if .Values.postgres.enabled -}}
postgres://{{ .Values.postgres.user }}:{{ .Values.postgres.password }}@{{ include "kirmya.postgresHost" . }}:{{ .Values.postgres.port }}/{{ .Values.postgres.db }}?sslmode=disable
{{- else -}}
{{- .Values.secrets.databaseUrl -}}
{{- end -}}
{{- end -}}
