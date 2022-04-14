# rut_wialon_gateway

installer
```shell
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/installer.sh && chmod +x installer.sh && ./installer.sh
rm -rf installer.sh /overlay/install
```

запуск
```shell
service r_w_g_service_controller enable
service r_w_g_service_controller start
```

geo data
```shell
gpsctl -h
```

компиляция для mips архитектуры
```shell
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -trimpath -ldflags="-s -w" 'cmd/main_app/main.go' && upx -9 main && mv main install/RWG_main_app
```