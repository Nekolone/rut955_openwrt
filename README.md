# rut955_openwrt

installer
```shell
wget https://raw.githubusercontent.com/Nekolone/rut955_openwrt/main/installer.sh && chmod +x installer.sh && ./installer.sh
service rut_wialon_gateway enable
service rut_wialon_gateway start
```

get gps data
```shell
cat /dev/ttyUSB2
```
or
```shell
gpsctl -h
```


* cmd/main_app/main.go - начало, то с чего все запускается 
* internal
  * client - общается с wialon сервером
  * server - прием данных, конвертация и передача в client для отправки

запуск на linux
```shell
go run cmd/main_app/main.go
```

компиляция для rut955
```shell
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -trimpath -ldflags="-s -w" 'cmd/main_app/main.go' && upx -9 main
```