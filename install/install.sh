#!/bin/sh

wget https://github.com/Nekolone/rut955_openwrt/raw/main/install/main_app
wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/install/module_custom_config.json
wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/install/module_mqtt_config.json
wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/install/rut_data_processing_service_config.json
wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/install/rut_gateway_config_paths.json
wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/install/rut_modules_config.json
wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/install/rut_wialon_client_config.json
wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/install/rut_wialon_gateway
wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/install/rut_wialon_gateway.sh


chmod +x rut_wialon_gateway
mv /overlay/install/rut_wialon_gateway /etc/init.d/rut_wialon_gateway
service rut_wialon_gateway enable

chmod +x main_app
mv /overlay/install/main_app /overlay/wialon_rut955_gateway/rut_wialon_gateway_app

mv /overlay/install/module_custom_config.json /overlay/wialon_rut955_gateway/module_custom_config.json
mv /overlay/installmodule_mqtt_config.json /overlay/wialon_rut955_gateway/module_mqtt_config.json
mv /overlay/installrut_data_processing_service_config.json /overlay/wialon_rut955_gateway/rut_data_processing_service_config.json
mv /overlay/installrut_gateway_config_paths.json /overlay/wialon_rut955_gateway/rut_gateway_config_paths.json
mv /overlay/installrut_modules_config.json /overlay/wialon_rut955_gateway/rut_modules_config.json
mv /overlay/installrut_wialon_client_config.json /overlay/wialon_rut955_gateway/rut_wialon_client_config.json

mv /overlay/installrut_wialon_gateway.sh /overlay/wialon_rut955_gateway/rut_wialon_gateway.sh


wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/curr_version
mv /overlay/install/curr_version /overlay/wialon_rut955_gateway/curr_version

service rut_wialon_gateway start

cd ..
#rm -rf install
