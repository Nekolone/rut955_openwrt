#!/bin/sh


/overlay/wialon_rut955_gateway/rut_wialon_gateway_app &


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
    pgrep "/overlay/wialon_rut955_gateway/rut_wialon_gateway_app" | xargs kill

    if ! wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/main_app ; then
        continue
    fi

    mv main_app rut_wialon_gateway_app
    chmod +x rut_wialon_gateway_app
    /overlay/wialon_rut955_gateway/rut_wialon_gateway_app &
done
