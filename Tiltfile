# vim: set filetype=bzl :

load('ext://restart_process', 'docker_build_with_restart')

CONTROLLER_DOCKERFILE = '''FROM golang:alpine
WORKDIR /app
COPY ./bin/service-unit .
CMD ["./service-unit"]
'''

# Generate manifests and go files
local_resource('service-unit manifests', cmd='make devmanifests', deps=['dev/config', 'bin/tbctl'])

# Deploy service units
watch_file('./dev/manifest/')
k8s_yaml(kustomize('./dev/manifest'))

# Re-compile go
local_resource('Watch & Compile', 'make devbuild', deps=['cmd', 'pkg', 'internal'])

# Build service-unit image
docker_build_with_restart(
    ref='hiroki11hanada/service-unit:dev',
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
    ref='hiroki11hanada/tb-load-generator:dev',
    context='./build/load-generator/',
)

# Build mongo image
docker_build(
    ref='hiroki11hanada/stateful-unit-mongo:dev',
    context='./build/stateful-unit/mongo/',
)
