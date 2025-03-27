#!/bin/bash

set -e

source $(dirname "$0")/_common.sh

KONG_IMAGE=${KONG_IMAGE_REPO:-kong/kong-gateway}:${KONG_IMAGE_TAG:-3.6}
NETWORK_NAME=kong-test

KONG_ROUTER_FLAVOR=${KONG_ROUTER_FLAVOR:-'traditional_compatible'}

PG_CONTAINER_NAME=pg
DATABASE_USER=kong
DATABASE_NAME=kong
KONG_DB_PASSWORD=kong
KONG_PG_HOST=pg

GATEWAY_CONTAINER_NAME=kong

KONG_ADMIN_API=http://localhost:8001

create_network

if [[ ! -z "${TEST_KONG_PULL_USERNAME}" ]]; then
  echo "${TEST_KONG_PULL_PASSWORD}" | docker login --username "${TEST_KONG_PULL_USERNAME}" --password-stdin
fi

function deploy_kong_ee()
{
  # Start Kong Gateway EE
  docker run -d --name $GATEWAY_CONTAINER_NAME \
    --network=$NETWORK_NAME \
    -e "KONG_DATABASE=postgres" \
    -e "KONG_PG_HOST=$KONG_PG_HOST" \
    -e "KONG_PG_USER=$DATABASE_USER" \
    -e "KONG_PG_PASSWORD=$KONG_DB_PASSWORD" \
    -e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
    -e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
    -e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
    -e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
    -e "KONG_ADMIN_LISTEN=0.0.0.0:8001" \
    -e "KONG_PORTAL_GUI_URI=127.0.0.1:8003" \
    -e "KONG_ADMIN_GUI_URL=http://127.0.0.1:8002" \
    -e "KONG_LICENSE_DATA=$KONG_LICENSE_DATA" \
    -e "KONG_ADMIN_GUI_AUTH=basic-auth" \
    -e "KONG_ENFORCE_RBAC=on" \
    -e "KONG_PORTAL=on" \
    -e "KONG_ADMIN_GUI_SESSION_CONF={}" \
    -e "KONG_ROUTER_FLAVOR=${KONG_ROUTER_FLAVOR}" \
    -p 8000:8000 \
    -p 8443:8443 \
    -p 8001:8001 \
    -p 8444:8444 \
    -p 8002:8002 \
    -p 8445:8445 \
    -p 8003:8003 \
    -p 8004:8004 \
    --label "$DOCKER_LABEL" \
    $KONG_IMAGE
}

deploy_pg
perform_migrations
deploy_kong_ee

waitContainer "Kong" 8001 0.2
waitKongAPI 0.5
