# GeoLite2-Go
GeoLite2 Integration Demo

# Общая структура проекта

## Методы сервиса

Сервис поддерживает два вида запросов:
1. `/geoip?ip={ip}` На API сервиса GeoLite2
2. `geoip_local?ip={ip}` К локальной БД GeoLite2 (БД должна быть предоставлена, автоматического получения нет)

Позитивный ответ имеет вид:
```json
{"ip":"2001:0:130f::9c0:876a:1301","country":"US","fresh":true}
```
- ip - полученный на вход ip адрес
- country - двухбуквенный код страны от сервиса GeoIP
- fresh - был ли ответ получен от сервиса (`true`) или из кэша (`false`)

Негативный ответ имеет вид:
```json
{"error":"параметр ip не является валидным ip адресом","errcode":400}
```
- error - текст ошибки
- errcode - код ошибки (http), который дублируется как http code в заголовке ответа

## Общее

Архитектурно всё максимально плоско так как нет никаких достаточно больших логически связных фрагментов которые стоит выделять в отдельные модули.

**Логи в исходниках** `log.Println` отправляют вывод в cli и как следствие отображаются в логах докера, что в целом удобно в любых обстоятельствах.


# Зависимости
Учитывая простоту проекта постарался свести зависимости к минимуму. Стандартный роутер из net/http более чем достаточен.

### Модуль `github.com/bradfitz/gomemcache`
Очевидных аналогов в репозитории Go не нашлось. Библиотека имеет встроенный mutex для синхронизации запросов и избегания race-conditions и в целом выглядит хорошо написанной. Не вижу смысла пытаться написать руками.

### Модуль `github.com/joho/godotenv`
Упрощает работу с .env файлами. В целом вариантов работы с env много: можно было подгружать напрямую в контейнер глобально, можно указать в Dockerfile или docker-compose.yml, можно было вручную читать .env файл. Дело привычки.

### Модуль `github.com/oschwald/geoip2-golang`
Выглядит как единственная опция для работы с файлами баз данных GeoIP2. Оно работает, разрабатывать свой Reader для mmdb файлов выглядит излишним на данном этапе.

### Остальное
Остальные модули являются непрямыми зависимостями от описанных выше.

# Контейнер

Это простой контейнер для разработки на базе полного контейнера Go + Debian Bullseye с поддержкой компиляции исходников. Рестарт контейнера пересобирает исходники посредством bash-скриптов (по этому используется версия с компилятором) и перезапускается утилитой `runit`. 

>Для деплоймента следует использовать runtime контейнер в который зашивает уже скомпилированный файл и зависимости.

Утилита runit позволяет запускать несколько сервисов и держать их более или менее живыми. По-верх стоит ещё делать health check на уровне Dockerfile, но опять же - это контейнер для разработки. 

## Конфигурация и запуск

Из файла `.env.example` в `./src/config` следует сформировать `.env` файл с настройками.

В .env файле можно изменить порт на котором находится сервис. В таком случае эту же правку надо внести в docker-compose.yml (проброс портов).

Сервис должен быть в состоянии работать с внешним Memcached, но я не проверял. Для настройки работы внутреннего memcached доступен файл /etc/memcached.conf

Всё остальное в целом не требует и/или не поддерживает настройку и работает как есть.

Для запуска потребуется docker compose. Запуск выполняется `docker compose up -d` или `docker compose up --build -d` если требуется пересобрать контейнер.