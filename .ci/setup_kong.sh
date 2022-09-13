#!/bin/bash

set -e
# download Kong deb

sudo apt-get update
sudo apt-get install openssl libpcre3 procps perl wget zlibc

function setup_kong(){
  SWITCH="1.3.000"
  SWITCH2="2.0.000"
  SWITCH3="2.8.000"

  URL="https://download.konghq.com/gateway-1.x-ubuntu-xenial/pool/all/k/kong/kong_${KONG_VERSION}_all.deb"

  if [[ "$KONG_VERSION" > "$SWITCH" ]];
  then
  URL="https://download.konghq.com/gateway-1.x-ubuntu-xenial/pool/all/k/kong/kong_${KONG_VERSION}_amd64.deb"
  fi

  if [[ "$KONG_VERSION" > "$SWITCH2" ]];
  then
  URL="https://download.konghq.com/gateway-2.x-ubuntu-xenial/pool/all/k/kong/kong_${KONG_VERSION}_amd64.deb"
  fi

  if [[ "$KONG_VERSION" > "$SWITCH3" ]];
  then
  URL="https://download.konghq.com/gateway-3.x-ubuntu-focal/pool/all/k/kong/kong_${KONG_VERSION}_amd64.deb"
  fi

  echo "Saving ${URL} to kong.deb"
  RESPONSE_CODE=$(/usr/bin/curl -sL \
    -w "%{http_code}" \
    $URL -o kong.deb)
  if [[ $RESPONSE_CODE != "200" ]]; then
    echo "error retrieving kong package from ${URL}. response code ${RESPONSE_CODE}"
    exit 1 
  fi
}

function setup_kong_enterprise(){
  KONG_VERSION="${KONG_VERSION#enterprise-}"
  SWITCH="1.5.0.100"
  SWITCH2="2.0.0.000"
  SWITCH3="2.8.0.000"

  URL="https://download.konghq.com/private/gateway-1.x-ubuntu-xenial/pool/all/k/kong-enterprise-edition/kong-enterprise-edition_${KONG_VERSION}_all.deb"

  if [[ "$KONG_VERSION" > "$SWITCH" ]];
  then
  URL="https://download.konghq.com/gateway-1.x-ubuntu-xenial/pool/all/k/kong-enterprise-edition/kong-enterprise-edition_${KONG_VERSION}_all.deb"
  fi

  if [[ "$KONG_VERSION" > "$SWITCH2" ]];
  then
  URL="https://download.konghq.com/gateway-2.x-ubuntu-xenial/pool/all/k/kong-enterprise-edition/kong-enterprise-edition_${KONG_VERSION}_all.deb"
  fi

  if [[ "$KONG_VERSION" > "$SWITCH3" ]];
  then
  URL="https://download.konghq.com/gateway-3.x-ubuntu-bionic/pool/all/k/kong-enterprise-edition/kong-enterprise-edition_${KONG_VERSION}_amd64.deb"
  fi

  echo "Saving ${URL} to kong.deb"
  RESPONSE_CODE=$(/usr/bin/curl -sL \
    -w "%{http_code}" \
    -u $KONG_ENTERPRISE_REPO_USERNAME:$KONG_ENTERPRISE_REPO_PASSSWORD \
    $URL -o kong.deb)
  if [[ $RESPONSE_CODE != "200" ]]; then
    echo "error retrieving kong enterprise package from ${URL}. response code ${RESPONSE_CODE}"
    exit 1 
  fi
}

if [[ $KONG_VERSION == *"enterprise"* ]]; then
  setup_kong_enterprise
else
  setup_kong
fi

sudo dpkg -i kong.deb
echo $KONG_LICENSE_DATA | sudo tee /etc/kong/license.json
export KONG_LICENSE_PATH=/tmp/license.json
export KONG_PASSWORD=kong
export KONG_ENFORCE_RBAC=on
export KONG_PORTAL=on

sudo KONG_PASSWORD=kong kong migrations bootstrap
sudo kong version
sudo KONG_ADMIN_GUI_SESSION_CONF='{}' KONG_ADMIN_GUI_AUTH=basic-auth KONG_ENFORCE_RBAC=on KONG_PORTAL=on kong start
