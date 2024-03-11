# Сервис SSO

___

## Структура проекта

```bash
.
├───cmd
│   ├───migrator
│   ├───proxy
│   └───sso
├───config
├───env
├───internal
│   ├───app
│   │   └───grpcapp
│   ├───config
│   ├───domain
│   │   └───models
│   ├───grpc
│   │   ├───handler
│   │   │   ├───auth
│   │   │   ├───permission
│   │   │   └───userInfo
│   │   └───interceptor
│   │       ├───auth
│   │       └───validation
│   ├───lib
│   │   ├───jwt
│   │   ├───logger
│   │   │   ├───handlers
│   │   │   │   ├───slogdiscard
│   │   │   │   └───slogpretty
│   │   │   └───sl
│   │   └───secret
│   ├───service
│   │   ├───auth
│   │   ├───permission
│   │   └───userInfo
│   └───storage
│       └───sqlite
├───migrations
└───storage
```

### Сервис предоставляет 7 эндпоитов

Можно делать как gRPC запросы (вызов метода), так и HTTP

Методы, что они принимают и что возвращают, можно посмотреть здесь: [интерфейс](https://github.com/dedmouze/protos)  
Протофайлы находятся [тут](https://github.com/dedmouze/protos/tree/main/proto/sso)