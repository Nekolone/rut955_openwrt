#!/bin/sh

chmod +x rut_wialon_gateway
mv rut_wialon_gateway /etc/init.d/rut_wialon_gateway
service rut_wialon_gateway enable

chmod +x main_app
mv main_app /overlay/wialon_rut955_gateway/rut_wialon_gateway_app

mv module_custom_config.json /overlay/wialon_rut955_gateway/module_custom_config.json
mv module_mqtt_config.json /overlay/wialon_rut955_gateway/module_mqtt_config.json
mv rut_data_processing_service_config.json /overlay/wialon_rut955_gateway/rut_data_processing_service_config.json
mv rut_gateway_config_paths.json /overlay/wialon_rut955_gateway/rut_gateway_config_paths.json
mv rut_modules_config.json /overlay/wialon_rut955_gateway/rut_modules_config.json
mv rut_wialon_client_config.json /overlay/wialon_rut955_gateway/rut_wialon_client_config.json

mv rut_wialon_gateway.sh /overlay/wialon_rut955_gateway/rut_wialon_gateway.sh



wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/curr_version
mv curr_version /overlay/wialon_rut955_gateway/curr_version

service rut_wialon_gateway start
