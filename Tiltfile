# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

### K8s Config ###

# Uncomment to use secrets
# k8s_yaml('./infra/development/k8s/secrets.yaml')

k8s_yaml('./infra/deploy/development/k8s/app-config.yaml')

### End of K8s Config ###
### API Gateway ###

gateway_compile_cmd = 'mkdir -p build && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/apigateway/cmd/main.go'

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/apigateway', './shared'], labels=["compiles"])


docker_build_with_restart(
  'ride-sharing/api-gateway',
  '.',
  entrypoint=['/app/build/api-gateway'],
  dockerfile='./infra/deploy/development/docker/api-gateway.Dockerfile',
  only=[
    './build/api-gateway',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/deploy/development/k8s/api-gateway-deployment.yaml')
k8s_resource('api-gateway', port_forwards=8081,
             resource_deps=['api-gateway-compile'], labels=["services"])
### End of API Gateway ###
### Trip Service ###

# Uncomment once we have a trip service

trip_compile_cmd = 'mkdir -p build && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/trip-service ./services/trip-service/cmd/main.go'

local_resource(
  'trip-service-compile',
  trip_compile_cmd,
  deps=['./services/trip-service', './shared'], labels=["compiles"])

docker_build_with_restart(
  'ride-sharing/trip-service',
  '.',
  entrypoint=['/app/build/trip-service'],
  dockerfile='./infra/deploy/development/docker/trip-service.Dockerfile',
  only=[
    './build/trip-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/deploy/development/k8s/trip-service-deployment.yaml')
k8s_resource('trip-service', resource_deps=['trip-service-compile'], labels=["services"])

### End of Trip Service ###
### Driver Service ###

driver_compile_cmd = 'mkdir -p build && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/driver-service ./services/driver-service/cmd'

local_resource(
  'driver-service-compile',
  driver_compile_cmd,
  deps=['./services/driver-service', './shared'], labels=["compiles"])

docker_build_with_restart(
  'ride-sharing/driver-service',
  '.',
  entrypoint=['/app/build/driver-service'],
  dockerfile='./infra/deploy/development/docker/driver-service.Dockerfile',
  only=[
    './build/driver-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/deploy/development/k8s/driver-service-deployment.yaml')
k8s_resource('driver-service', resource_deps=['driver-service-compile'], labels=["services"])

### End of Driver Service ###
### Web Frontend ###

docker_build_with_restart(
  'ride-sharing/web',
  './web',
  entrypoint=['npm', 'run', 'dev', '--', '--hostname', '0.0.0.0'],
  dockerfile='./infra/deploy/development/docker/web.Dockerfile',
  live_update=[
    sync('./web/src', '/app/src'),
    sync('./web/public', '/app/public'),
    sync('./web/package.json', '/app/package.json'),
    sync('./web/package-lock.json', '/app/package-lock.json'),
    run('npm install', trigger=['./web/package.json', './web/package-lock.json']),
  ],
)

k8s_yaml('./infra/deploy/development/k8s/web-deployment.yaml')
k8s_resource('web', port_forwards=3000, labels=["frontend"])

### End of Web Frontend ###
