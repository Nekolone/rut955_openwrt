# rut955_openwrt

installer
```shell
wget https://raw.githubusercontent.com/Nekolone/rut_wialon_gateway/main/installer.sh && chmod +x installer.sh && ./installer.sh
service r_w_g_service_controller enable
service r_w_g_service_controller start
rm -rf installer.sh /overlay/install
```

```shell
gpsctl -h
```

компиляция для rut955
```shell
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -trimpath -ldflags="-s -w" 'cmd/main_app/main.go' && upx -9 main
```