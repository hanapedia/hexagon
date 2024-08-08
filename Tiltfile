# vim: set filetype=bzl :

load('ext://restart_process', 'docker_build_with_restart')

# Define default values for your arguments
config.define_string('docker_user', usage='Docker user name used to build images for local testing.')

# Parse the arguments
cfg = config.parse()

# Access the arguments
docker_user = cfg.get('docker_user', 'hexagonbenchmark')

CONTROLLER_DOCKERFILE = '''FROM golang:alpine
WORKDIR /app
COPY ./bin/service-unit .
CMD ["./service-unit"]
'''

DATAGEN_DOCKERFILE = '''FROM golang:1.22
WORKDIR /app
COPY ./bin/datagen .
'''

# Generate manifests and go files
local_resource('service-unit manifests', cmd='make devmanifests DOCKER_USER=%s' % docker_user, deps=['dev/config', 'bin/hexctl'])

# Deploy service units
watch_file('./dev/manifest/')
k8s_yaml(kustomize('./dev/manifest'))

# Re-compile service unit
local_resource('Watch & Compile service-unit', 'make devbuild', deps=['cmd', 'pkg', 'internal'])

# Re-compile hexctl
local_resource('Watch & Compile hexctl', 'make devbuildcli', deps=['cmd/hexctl', 'pkg/hexctl', 'pkg/operator', 'internal/hexctl'])

# Re-compile datage
local_resource('Watch & Compile datagen', 'make devbuilddatagen', deps=['cmd/datagen', 'internal/datagen'])

# Build service-unit image
docker_build_with_restart(
    ref='%s/service-unit:dev' % docker_user,
    context='.',
    dockerfile_contents=CONTROLLER_DOCKERFILE,
    entrypoint=['./service-unit'],
    only=['./bin/service-unit'],
    live_update=[
        sync('./bin/service-unit', '/app/service-unit'),
    ]
)

# Build load-generator image
docker_build(
    ref='%s/tb-load-generator:dev' % docker_user,
    context='./build/load-generator/',
)

# Build datagen image
docker_build(
    ref='%s/datagen:dev' % docker_user,
    context='.',
    dockerfile_contents=DATAGEN_DOCKERFILE,
    only=['./bin/datagen'],
    live_update=[
        sync('./bin/datagen', '/app/datagen'),
    ]
)

# Build mongo image
docker_build(
    ref='%s/stateful-unit-mongo:dev' % docker_user,
    context='./build/stateful-unit/mongo/',
    build_args={
        'BUILDER_IMAGE': 'datagen:dev'
    },
)

# Build redis image
docker_build(
    ref='%s/stateful-unit-redis:dev' % docker_user,
    context='./build/stateful-unit/redis/',
    build_args={
        'BUILDER_IMAGE': 'datagen:dev'
    },
)
