# rut955_openwrt

installer
```shell
wget https://github.com/Nekolone/rut955_openwrt/archive/master.tar.gz && tar -xf master.tar.gz && mv rut955_openwrt-main/install install && rm -rf master.tar.gz rut955_openwrt-main
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