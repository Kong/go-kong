#!/bin/bash

# Usage: waitContainer "PostgreSQL" 5432 0.2
function waitContainer()
{
  local container=${1}
  local port=${2}
  local sleep_time=${3}

  for try in {1..100}; do
    echo "waiting for ${container}.."
    nc localhost ${port} && break;
    sleep ${sleep_time}
  done
}

function create_network()
{
  # create docker network if it doesn't exist
  if [[ -z $(docker network ls --filter name=${NETWORK_NAME} -q) ]]; then
    docker network create ${NETWORK_NAME}
  fi
}

function waitKongAPI() {
  for try in {1..100}; do
    echo "waiting for Kong Admin API.."
    curl -f -LI $KONG_ADMIN_API -H kong-admin-token:$KONG_ADMIN_TOKEN && break;
    sleep $1
  done
}

function deploy_pg()
{
  # Start a PostgreSQL container
  docker run --rm -d --name $PG_CONTAINER_NAME \
    --network=$NETWORK_NAME \
    -p 5432:5432 \
    -e "POSTGRES_USER=$DATABASE_USER" \
    -e "POSTGRES_DB=$DATABASE_NAME" \
    -e "POSTGRES_PASSWORD=$KONG_DB_PASSWORD" \
    postgres:9.6

  waitContainer "PostgreSQL" 5432 0.2
}

function perform_migrations()
{
  for try in {1..10}; do
    # Prepare the Kong database
    docker run --rm --network=$NETWORK_NAME \
      -e "KONG_DATABASE=postgres" \
      -e "KONG_PG_HOST=$KONG_PG_HOST" \
      -e "KONG_PG_PASSWORD=$KONG_DB_PASSWORD" \
      -e "KONG_PASSWORD=$KONG_DB_PASSWORD" \
      -e "KONG_LICENSE_DATA=$KONG_LICENSE_DATA" \
      $KONG_IMAGE kong migrations bootstrap && break
  done
}
