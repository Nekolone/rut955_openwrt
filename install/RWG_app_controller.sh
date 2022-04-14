#!/bin/sh

#TODO: сделать запуск если нет процесса RWG_main_app


/overlay/rut_wialon_gateway/RWG_main_app &
/overlay/rut_wialon_gateway/RWG_buff_remover.sh &

app_path=/overlay/rut_wialon_gateway
buf_path=/tmp/RWG_app_buffer
rep_path=https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway

while :; do

  sleep 60m
  rm -rf last_version

  if ! wget $rep_path/main/last_version; then
    continue
  fi

  UPD_VER=$(cat last_version)
  CUR_VER=$(cat $app_path/curr_version)

  if [[ $UPD_VER == $CUR_VER ]]; then
    continue
  fi

  echo "$UPD_VER" >$app_path/curr_version
  pgrep "RWG_main_app" | xargs kill

  mv $app_path/RWG_main_app $buf_path/RWG_main_app

  if ! wget $rep_path/main/RWG_main_app; then
    mv $buf_path/RWG_main_app $app_path/RWG_main_app
    continue
  fi

  mv last_version $app_path/curr_version

  chmod +x RWG_main_app
  mv RWG_main_app $app_path/RWG_main_app

  /overlay/rut_wialon_gateway/RWG_main_app &

  rm $app_path/RWG_main_app

done
