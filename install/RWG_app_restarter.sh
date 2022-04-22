#!/bin/sh

while :; do
  if [[ -z $(pgrep -f RWG_main_app) ]] ; then
    /overlay/rut_wialon_gateway/RWG_main_app &
  fi

  sleep 5m
done
