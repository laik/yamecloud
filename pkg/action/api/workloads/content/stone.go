package content

const (
	stoneTpl = `kind: Stone
apiVersion: nuwa.nip.io/v1
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.Name}}
    app-uuid: {{.UUID}}
    yce-cloud-extensions: {{.CDName}}
spec:
  template:
    metadata:
      name: {{.Name}}
      labels:
        app: {{.Name}}
        app-uuid: {{.UUID}}
    spec:
      containers:
        - name: {{.Name}}
          image: {{.Image}}
          {{- if .Commands}}
          command:
          {{ range .Commands}}
            - {{.}}
          {{ end }}
          {{- end }}
          {{- if .Args}}
          args:
          {{ range .Args}}
            - {{.}}
          {{ end }}
          {{- end }}
          {{- if .Environments}}
          env:
          {{range .Environments}}
            - name: {{.Name}}
              value: {{.Envvalue}}
          {{ end }}
          {{- end }}
          resources:
            limits:
              cpu: {{.CpuLimit}}
              memory: {{.MemoryLimit}}
            requests:
              cpu: {{.CpuRequests}}
              memory: {{.MemoryRequests}}
          imagePullPolicy: Always
          {{- if .ConfigVolumes}}
          volumeMounts:
            {{range .ConfigVolumes}}
            - name: {{.MountName}}
              mountPath: {{.MountPath}}
              subPath: {{.SubPath}}
            {{ end }}
          {{- end }}
      {{- if .ConfigVolumes}}
      volumes:
        {{range .ConfigVolumes}}
        - name: {{.MountName}}
          configMap:
            name: {{$.Name}}
            {{- if .CMItems}}
            items:
              {{range .CMItems}}
              - key: {{.VolumeName}}
                path: {{.VolumePath}}
              {{ end }}
            {{- end }}
        {{ end }}
      {{- end }}
  strategy: Release
  coordinates:
{{range .Coordinates}}
    - group: {{.Group}}
      zoneset:
{{range .ZoneSet}}
        - zone: {{.Zone}}
          rack: {{.Rack}}
          host: {{.Host}}
{{end}}
      replicas: {{.Replicas}}
{{end}}
  service:
    ports:
{{range .ServicePorts}}
      - name: {{.Name}}
        protocol: {{.Protocol}}
        port: {{.Port}}
        targetPort: {{.TargetPort}}
{{end}}
    type: {{.ServiceType}}`
)
