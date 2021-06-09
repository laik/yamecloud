package content

const (
	StoneTpl = `
kind: Stone
apiVersion: nuwa.nip.io/v1
metadata:
  name: {{.name}}
  namespace: {{.namespace}}
  labels:
    app: {{.name}}
    uuid: {{.uuid}}
spec:
  strategy: {{or .strategy "Alpha"}}

  template:
    metadata:
      name: {{.name}}
      labels:
        app: {{.name}}
        uuid: {{.uuid}}

    spec:
      {{- if .containers}}
      containers:
        {{range $index,$container := .containers}}
        {{$container_name := (printf "container-%v" $index)}}
        - name: {{or .name $container_name}}
          image: {{.image}}
          imagePullPolicy: {{or .image_pull_policy "Always"}}

# container commands
          {{- if .commands}}
          command:
            {{ range .commands}}
            - {{.}}
            {{end}}
          {{- end }}

# container args
          {{- if .args}}
          args:
            {{ range .args}}
            - {{.}}
            {{end}}
          {{- end }}

#container environments
          {{- if .environments}}
          env:
          {{range .environments}}
          {{$_value := (printf "'%v'" .value)}}
            - name: {{.name}}
              value: {{$_value}}
          {{end}}
          {{- end }}

#container resource limits
          {{- if .resource_limits}}
          {{$cpu_limit := (printf "%vm" .resource_limits.cpu_limit)}}
          {{$memory_limit := (printf "%vM" .resource_limits.memory_limit)}}
          {{$cpu_request := (printf "%vm" .resource_limits.cpu_request)}}
          {{$memory_request := (printf "%vM" .resource_limits.memory_request)}}
          resources:
            limits:
              cpu: {{$cpu_limit}}
              memory: {{$memory_limit}}
            requests:
              cpu: {{$cpu_request}}
              memory: {{$memory_request}}
          {{- end}}

#container volumes mount
          {{- if .volume_mounts}}
          volumeMounts:
            {{range $index,$value := .volume_mounts}}
            - mountPath: {{.mount_path}}
              name: {{.name}}
              {{- if .sub_path}}
              subPath: {{.sub_path}}
              {{- end}}
            {{end}}
          {{- end}}

# containers range end
        {{end}}
# if containers end
      {{- end}}

# pull secrets
      {{- if .image_pull_secrets}}
      imagePullSecrets:
        {{range .image_pull_secrets}}
        - name: {{.}}
        {{end}}
      {{- end}}


# volumes
      {{- if .volumes}}
      volumes:
        {{range .volumes}}
        - name: {{.name}}

# configmap or other volumes
          {{- if .configmap}}
          configMap:
            name: {{.configmap.name}}
          {{- if .configmap.items}}
            items:
            - key: {{.configmap.items.key}}
              path: {{.configmap.items.path}}
          {{- end}}
          {{- end}}

# secret 
          {{- if .secret}}
          secret:
            name: {{.secret.name}}
          {{- if .secret.items}}
            items:
            - key: {{.secret.items.key}}
              path: {{.secret.items.path}}
          {{- end}}
          {{- end}}

        {{end}}
      {{- end }}


# volumeClaimTemplates
{{- if .volume_claim_templates}}
  volumeClaimTemplates:
    {{range .volume_claim_templates}}
    {{$size := (printf "%sMi" .size)}}
    - metadata:
        name: {{.metadata_name}}
      spec:
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: {{$size}}
        storageClassName: {{.storage_class_name}}
    {{end}}
{{- end}}

# containers
{{- if .coordinates}}
  coordinates:
  {{range .coordinates}}
    - group: {{.group}}
      zoneset:
      {{range .zoneset}}
        - zone: {{.zone}}
          rack: {{.rack}}
          host: {{.host}}
      {{end}}
      replicas: {{or .replicas 0}}
  {{end}}
{{end}}

# service spec
{{- if .service_ports}}
  service:
    ports:
      {{range $index,$value := .service_ports}}
      {{$port_name := (printf "port-%v" $index)}}
      - name: {{or .name $port_name }}
        protocol: {{or .protocol "TCP"}}
        port: {{.port}}
        targetPort: {{.target_port}}
      {{end}}
    type: {{or .service_type "ClusterIP"}}
{{end}}
`

	DeploymentTpl = `
kind: Deployment
apiVersion: apps/v1
metadata:
  name: {{.name}}
  namespace: {{.namespace}}
  labels:
    app: {{.name}}
spec:
  replicas: {{or (index .coordinates 0).replicas 1}}
  selector:
    matchLabels:
      app: {{.name}}
  template:
    metadata:
      name: {{.name}}
      labels:
        app: {{.name}}
    spec:
      {{- if .containers}}
      containers:
        {{range $index,$container := .containers}}
        {{$container_name := (printf "container-%v" $index)}}
        - name: {{or .name $container_name}}
          image: {{.image}}
          imagePullPolicy: {{or .image_pull_policy "Always"}}

# container port
          {{- if $.service_ports}}
          ports:
            {{range $.service_ports}}
            - containerPort: {{.target_port}}
              protocol: {{or .protocol "TCP"}}
            {{end}}
          {{- end }}

# container args
          {{- if .args}}
          args:
            {{ range .args}}
            - {{.}}
            {{end}}
          {{- end }}

#container environments
          {{- if .environments}}
          env:
          {{range .environments}}
          {{$_value := (printf "'%v'" .value)}}
            - name: {{.name}}
              value: {{$_value}}
          {{end}}
          {{- end }}

#container volumes mount
          {{- if .volume_mounts}}
          volumeMounts:
            {{range $index,$value := .volume_mounts}}
            - mountPath: {{.mount_path}}
              name: {{.name}}
              {{- if .sub_path}}
              subPath: {{.sub_path}}
              {{- end}}
            {{end}}
          {{- end}}

# containers range end
        {{end}}
# if containers end
      {{- end}}

# volumes
      {{- if .volumes}}
      volumes:
        {{range .volumes}}
        - name: {{.name}}

# configmap or other volumes
          {{- if .configmap}}
          configMap:
            name: {{.configmap.name}}
          {{- if .configmap.items}}
            items:
            - key: {{.configmap.items.key}}
              path: {{.configmap.items.path}}
          {{- end}}
          {{- end}}

# secret 
          {{- if .secret}}
          secret:
            name: {{.secret.name}}
          {{- if .secret.items}}
            items:
            - key: {{.secret.items.key}}
              path: {{.secret.items.path}}
          {{- end}}
          {{- end}}
        {{end}}
      {{- end }}
`

	ServiceTpl = `
kind: Service
apiVersion: v1
metadata:
  name: {{.name}}
  namespace: {{.namespace}}
spec:
# service spec
{{- if .service_ports}}
    ports:
      {{range $index,$value := .service_ports}}
      {{$port_name := (printf "port-%v" $index)}}
      - name: {{or .name $port_name }}
        protocol: {{or .protocol "TCP"}}
        port: {{.port}}
        targetPort: {{.target_port}}
      {{end}}
    type: {{or .service_type "ClusterIP"}}
{{end}}
    selector:
      app: {{.name}}
`
)
