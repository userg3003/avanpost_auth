# avanpost_auth
Авторизация через Avanpost FAM

## Запуск сервиса

### Конфигурация

Параметры конфигурации сервиса могут быть заданы как через конфигурационный файл,
так и посресдством переменных окружения.

**SERVICE_SHEMA** - схема обращения к сервису   
**SERVICE_HOST**  - хост  сервиса
**SERVICE_PORT**  - порт  сервиса
**SERVICE_OAUTH2_REDIRECT** - url для передачи ответа от Avanpost FAM  
**SERVICE_COOKIE_SESSION_NAME** - имя сессионной cookie для сохранения данных авторизации    
**SERVICE_COOKIE_SESSION_SECRET** - секрет для   сессионной cookie  
**SERVICE_REDIRECT_URL_AFTER_AUTH** - путь для перенаправления после успешной авторизации    
**OAUTH2_URL_AUTH_SHEMA** -  схема обращения к Avanpost FAM  
**OAUTH2_URL_AUTH_HOST** - хост Avanpost FAM  
**OAUTH2_URL_AUTH_PORT** - порт Avanpost FAM  
**OAUTH2_URL_AUTH_PATH** - путь для запроса авторизации в Avanpost FAM     
**OAUTH2_URL_TOKEN_PATH** - путь для в Avanpost FAM запроса токена       
**OAUTH2_URL_INFO_PATH** - путь для в Avanpost FAM запроса данных авторизации по токену  
**OAUTH2_CLIENT_ID** -   id клиента  в Avanpost FAM  
**OAUTH2_CLIENT_SECRET**   -   секрет клиента  в Avanpost FAM     
**SWAGGER** -   включить/отключить swagger (true/false)  
**LOG_LEVEL** - уроввень логирования  



## Команды make

***make build-service*** - собрать сервис   
***make build-oauth2-serever*** - собрать тестовый OAuth2 сервер  
***make run-service*** - запустить сервис  
***make run-oauth2-server*** - запустить    тестовый OAuth2 сервер  
***make swag*** - сгенерировать документацию swagger       


### Запуск сервиса
> make run-service  

Для тестирования сервиса следует запустить тестовый OAuth2-сервер.   
> make run-oauth2-server

В сервере захардкожены: 
 - юзер/пароль (test/test)
 - id клиента (OAUTH2_CLIENT_ID=1234)
 - url для передачи ответа от Avanpost FAM (appauth)
 - секрет (OAUTH2_CLIENT_SECRET=aabbccdd)  
 - *OAUTH2_URL_INFO_PATH* (info)
 - порт на котором запущен сервер (14000)


## Генерация swagger

Для генерации swagger из исхоодного кода необходим необходим
конвертер аннотаций Go в документацию swagger [swag](https://github.com/swaggo/swag).  
Его можно установить так:

> $ go get -u github.com/swaggo/swag/cmd/swag

или для go начиная с 1.16

> $ go install github.com/swaggo/swag/cmd/swag@latest

или скачать скомпилированный вариант [отсюда](https://github.com/swaggo/swag/releases).


#### Генерация докуметации

> make swag


После генерации документации Swagger будет доступен  по адресу: [http://***host:port***/swagger/index.html](http://<host>:<port>/swagger/index.html)

 