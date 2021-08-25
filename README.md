# Sysmon - демон, осуществляющий "Системный мониторинг"

Демон - программа, собирающая информацию о системе, на которой запущена, и отправляющая её своим клиентам по GRPC.

![go workflow](https://github.com/artofey/sysmon/actions/workflows/go.yml/badge.svg)
[![GoDoc](https://godoc.org/github.com/artofey/sysmon?status.svg)](https://pkg.go.dev/github.com/artofey/sysmon)

## Описание
Клиент после подключения должен указать два параметра:
- раз во сколько секунд он хочет получать информацию (N)
- усредненной за сколько секунд должна быть информация (M)

> Например, N = 5с, а M = 15с, тогда демон "молчит" первые 15 секунд, затем выдает снапшот за 0-15с; через 5с (в 20с) выдает снапшот за 5-20с; через 5с (в 25с) выдает снапшот за 10-25с и т.д.

Что может собирать:
- [x] Средняя загрузка системы (load average).

- [x] Средняя загрузка CPU (%user_mode, %system_mode, %idle).

- [ ] Загрузка дисков:
    - tps (transfers per second);
    - KB/s (kilobytes (read+write) per second);
    - CPU (%user_mode, %system_mode, %idle).

- [ ] Информация о дисках по каждой файловой системе:
    - использовано мегабайт, % от доступного количества;
    - использовано inode, % от доступного количества.

- [ ] Top talkers по сети:
    - по протоколам: protocol (TCP, UDP, ICMP, etc), bytes, % от sum(bytes) за последние **M**), сортируем по убыванию процента;
    - по трафику: source ip:port, destination ip:port, protocol, bytes per second (bps), сортируем по убыванию bps.

- [ ] Статистика по сетевым соединениям:
    - слушающие TCP & UDP сокеты: command, pid, user, protocol, port;
    - количество TCP соединений, находящихся в разных состояниях (ESTAB, FIN_WAIT, SYN_RCV и пр.).

Статистика представляет собой объекты, описанные в формате Protobuf.

## Поддерживаемая ОС
- [x] Linux.
- [x] Windows.
- [ ] Darwin.

## Запуск через docker-compose

```
docker-compose up --build
docker exec -it sysmon-client sh
./sysmon-client

```

## Установка сервера
Установка через **go get**:
```
go get -u github.com/artofey/sysmon/cmd/sysmon-server
```
Установка через **Docker**:
```
git clone https://github.com/artofey/sysmon.git
cd sysmon
make run-img-server
```
## Установка клиента
Установка через **go get**:
```
go get -u github.com/artofey/sysmon/cmd/sysmon-client
```
Установка через **Docker**:
```
git clone https://github.com/artofey/sysmon.git
cd sysmon
make run-img-client
```

<!-- ## Конфигурация
- Через аргументы командной строки можно указать, на каком порту стартует сервер.
- Через файл можно указать, какие из подсистем сбора включены/выключены. -->

## DEV requirements
- **protoc** `sudo apt install protobuf-compiler` ([other instalation options](https://grpc.io/docs/protoc-installation/))
