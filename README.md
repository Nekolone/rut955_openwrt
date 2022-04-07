# rut955_openwrt

installer
```shell
wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/installer.sh && chmod +x installer.sh && ./installer.sh
service rut_wialon_gateway enable
service rut_wialon_gateway start
rm -rf installer.sh /overlay/install
```

```shell
gpsctl -h
```

компиляция для rut955
```shell
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -trimpath -ldflags="-s -w" 'cmd/main_app/main.go' && upx -9 main
```