#!/bin/sh


./rut_wialon_gateway &


while :
do
    sleep 60m
    rm -rf last_version

    if ! wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/last_version ; then
        continue
    fi

    UPD_VER=$(cat last_version)
    CUR_VER=$(cat curr_version)

    if [[ $UPD_VER == $CUR_VER ]]; then
        continue
    fi

    echo "$UPD_VER" > /overlay/wialon_rut955_gateway/curr_version
    pgrep "./rut_wialon_gateway" | xargs kill

    if ! wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/main_app ; then
        continue
    fi

    mv main_app rut_wialon_gateway
    chmod +x rut_wialon_gateway
    ./rut_wialon_gateway &
done
