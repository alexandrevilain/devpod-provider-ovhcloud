name: ovhcloud
version: {{ .Env.VERSION }}
description: |-
  DevPod on OVHCloud
icon: https://avatars3.githubusercontent.com/ovh
optionGroups:
  - options:
      - AGENT_PATH
      - INACTIVITY_TIMEOUT
      - INJECT_DOCKER_CREDENTIALS
      - INJECT_GIT_CREDENTIALS
    name: "Agent options"
    defaultVisible: false
  - options:
    - OVHCLOUD_ENDPOINT
    - OVHCLOUD_APP_KEY
    - OVHCLOUD_APP_SECRET
    - OVHCLOUD_CONSUMER_KEY
    - OVHCLOUD_SERVICE_NAME
    - OVHCLOUD_REGION
    - OVHCLOUD_FLAVOR
    name: "OVHCloud options"
    defaultVisible: true
options:
  OVHCLOUD_ENDPOINT:
    description: The ovhcloud endpoint to use
    required: true
    password: false
    suggestions:
      - ovh-eu
      - ovh-us
      - ovh-ca
    default: ovh-eu
  OVHCLOUD_APP_KEY:
    description: The ovhcloud application key
    required: true
    password: true
    default: ""
  OVHCLOUD_APP_SECRET:
    description: The ovhcloud application secret
    required: true
    password: true
    default: ""
  OVHCLOUD_CONSUMER_KEY:
    description: The ovhcloud consumer key
    required: true
    password: true
    default: ""
  OVHCLOUD_SERVICE_NAME:
    description: The ovhcloud service name
    required: true
    password: false
    default: ""
  OVHCLOUD_REGION:
    description: The ovhcloud region to create the VM in. E.g. GRA11
    required: true
    default: ""
    suggestions:
    - SBG1
    - SBG3
    - SBG5
    - SBG7
    - GRA1
    - GRA3
    - GRA5
    - GRA7
    - GRA9
    - GRA11
    - UK1
    - DE1
    - WAW1 
    - BHS1
    - BHS2
    - BHS3
    - BHS5
    - VIN1
    - HIL1 
    - SGP1
    - SYD1
  OVHCLOUD_FLAVOR:
    description: The machine type to use.
    default: b2-7
    suggestions:
    - b2-120
    - b2-120-flex
    - b2-15
    - b2-15-flex
    - b2-30
    - b2-30-flex
    - b2-60
    - b2-60-flex
    - b2-7
    - b2-7-flex
    - bm-l1
    - bm-m1
    - bm-s1
    - c2-120
    - c2-120-flex
    - c2-15
    - c2-15-flex
    - c2-30
    - c2-30-flex
    - c2-60
    - c2-60-flex
    - c2-7
    - c2-7-flex
    - d2-2
    - d2-4
    - d2-8
    - g1-15
    - g1-30
    - g2-15
    - g2-30
    - g3-120
    - g3-30
    - i1-180
    - i1-45
    - i1-90
    - r2-120
    - r2-120-flex
    - r2-15
    - r2-15-flex
    - r2-240
    - r2-240-flex
    - r2-30
    - r2-30-flex
    - r2-60
    - r2-60-flex
    - s1-2
    - s1-4
    - s1-8
    - t1-le-45
    - t1-le-90
    - t1-le-180
    - t1-45
    - t1-90
    - t1-180
    - t2-le-45
    - t2-le-90
    - t2-le-180
    - t2-45
    - t2-90
    - t2-180
    - a100-180
    - a100-360
  INACTIVITY_TIMEOUT:
    description: If defined, will automatically stop the VM after the inactivity period.
    default: 10m
  INJECT_GIT_CREDENTIALS:
    description: "If DevPod should inject git credentials into the remote host."
    default: "true"
  INJECT_DOCKER_CREDENTIALS:
    description: "If DevPod should inject docker credentials into the remote host."
    default: "true"
  AGENT_PATH:
    description: The path where to inject the DevPod agent to.
    default: /var/lib/toolbox/devpod
agent:
  path: ${AGENT_PATH}
  dataPath: /home/debian/.devpod
  inactivityTimeout: ${INACTIVITY_TIMEOUT}
  injectGitCredentials: ${INJECT_GIT_CREDENTIALS}
  injectDockerCredentials: ${INJECT_DOCKER_CREDENTIALS}
  binaries:
    OVHCLOUD_PROVIDER:
{{- range file.Walk "./dist" -}}
{{- if not (file.IsDir .) -}}
    {{- $parts := . | regexp.Split "_" -1 }}
    {{- if eq "linux" (index $parts 1) }}
    - os: {{ index $parts 1 }}
      arch: {{ index $parts 2 }}
      path: https://github.com/alexandrevilain/devpod-provider-ovhcloud/releases/download/{{ $.Env.VERSION }}/{{ . | filepath.Base }}
      checksum: {{ . | file.Read | crypto.SHA256 }}
    {{- end -}}
{{- end -}}
{{- end }}
  exec:
    shutdown: |-
      ${OVHCLOUD_PROVIDER} stop
binaries:
  OVHCLOUD_PROVIDER:
{{- range file.Walk "./dist" -}}
{{- if not (file.IsDir .) -}}
    {{- $parts := . | regexp.Split "_" -1 }}
    {{- $ext := filepath.Ext . }}
    - os: {{ index $parts 1 }}
      arch: {{ index $parts 2 | strings.Trim $ext }}
      path: https://github.com/alexandrevilain/devpod-provider-ovhcloud/releases/download/{{ $.Env.VERSION }}/{{ . | filepath.Base }}
      checksum: {{ . | file.Read | crypto.SHA256 }}
{{- end -}}
{{- end }}
exec:
  init: ${OVHCLOUD_PROVIDER} init
  command: ${OVHCLOUD_PROVIDER} command
  create: ${OVHCLOUD_PROVIDER} create
  delete: ${OVHCLOUD_PROVIDER} delete
  start: ${OVHCLOUD_PROVIDER} start
  stop: ${OVHCLOUD_PROVIDER} stop
  status: ${OVHCLOUD_PROVIDER} status

  
