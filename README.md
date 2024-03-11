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
│   │   └───logger
│   │       ├───handlers
│   │       │   ├───slogdiscard
│   │       │   └───slogpretty
│   │       └───sl
│   ├───service
│   │   ├───auth
│   │   ├───permission
│   │   └───userInfo
│   └───storage
│       └───sqlite
├───migrations
├───storage
└───tests
    ├───migrations
    └───suite
```

### Сервис предоставляет 6 эндпоитов

Можно делать как gRPC запросы (вызов метода), так и HTTP

Методы, что они принимают и что возвращают, можно посмотреть здесь: [интерфейс](https://github.com/dedmouze/protos)  
Протофайлы находятся [тут](https://github.com/dedmouze/protos/tree/main/proto/sso)