#!/bin/sh

mkdir /overlay/rut_wialon_gateway
mkdir /tmp/RWG_app_buffer

wget https://github.com/Nekolone/rut_wialon_gateway/raw/main/install/RWG_main_app
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/install/MODULE_custom.json
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/install/MODULE_MQTT.json
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/install/CFG_data_processing_service.json
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/install/APP_PATHS.json
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/install/MODULES_LIST.json
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/install/CFG_wilaon_client.json
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/install/r_w_g_service_controller
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/install/RWG_app_controller.sh
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/install/RWG_buff_remover.sh
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/curr_version

mv ./MODULE_custom.json /overlay/rut_wialon_gateway/MODULE_custom.json
mv ./MODULE_MQTT.json /overlay/rut_wialon_gateway/MODULE_MQTT.json
mv ./CFG_data_processing_service.json /overlay/rut_wialon_gateway/CFG_data_processing_service.json
mv ./APP_PATHS.json /overlay/rut_wialon_gateway/APP_PATHS.json
mv ./MODULES_LIST.json /overlay/rut_wialon_gateway/MODULES_LIST.json
mv ./CFG_wilaon_client.json /overlay/rut_wialon_gateway/CFG_wilaon_client.json

chmod +x RWG_main_app
mv ./RWG_main_app /overlay/rut_wialon_gateway/RWG_main_app

chmod +x r_w_g_service_controller
mv ./r_w_g_service_controller /etc/init.d/r_w_g_service_controller

chmod +x RWG_buff_remover.sh
mv ./RWG_buff_remover.sh /overlay/rut_wialon_gateway/RWG_buff_remover.sh

chmod +x RWG_app_controller.sh
mv ./RWG_app_controller.sh /overlay/rut_wialon_gateway/RWG_app_controller.sh

mv ./curr_version /overlay/rut_wialon_gateway/curr_version
