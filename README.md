# Como usar

## Download dos módulos para cache local

```shell
make install
```

## Inicialização

```shell
make start
``` 

## Execução de testes

```shell
make run-test
``` 

# Curls

## CRUD de documentos

### Cadastro de documento

```shell
curl --location 'http://localhost:8080/api/document' \
--header 'Content-Type: application/json' \
--data '{
    "type": "CPF",
    "value": "62546948083",
    "isBlocked": false
}'
```

#### Resposta

```json
{
    "id": "874f1cbc-3333-433a-bdd9-06bf5231084f",
    "type": "CPF",
    "value": "62546948083",
    "createdAt": "2024-07-19T20:18:37.332934625Z",
    "updatedAt": "2024-07-19T20:18:37.332934744Z",
    "isBlocked": false
}
```

### Listagem de documentos

```shell
curl --location 'http://localhost:8080/api/document' \
--header 'Content-Type: application/json' \
--data ''
```

```shell
curl --location 'http://localhost:8080/api/document?id=5e337efe-e7ad-46f6-9ac9-71e8c2cb9634' \
--header 'Content-Type: application/json' \
--data ''
```

```shell
curl --location 'http://localhost:8080/api/document?value=99272744000122' \
--header 'Content-Type: application/json' \
--data ''
```

```shell
curl --location 'http://localhost:8080/api/document?value=0001&size=5&page=2' \
--header 'X-Timezone: America/Sao_Paulo'
```

#### Resposta

```json
[
    {
        "id": "3f612199-5a3a-412b-89cb-3110b786d84d",
        "type": "CNPJ",
        "value": "41.754.227/0001-00",
        "createdAt": "19/07/2024 16:17:03",
        "updatedAt": "19/07/2024 16:17:03",
        "isBlocked": false
    }
]
```

### Atualização de documento

```shell
curl --location --request PUT 'http://localhost:8080/api/document' \
--header 'Content-Type: application/json' \
--data '{
    "id": "106b8780-3bd7-4e73-a7df-726120ffc653",
    "value": "05.263.076/0001-23",
    "isBlocked": true
}'
```

#### Resposta

```json
{
    "id": "106b8780-3bd7-4e73-a7df-726120ffc653",
    "value": "05263076000123",
    "createdAt": "0001-01-01T00:00:00Z",
    "updatedAt": "2024-07-19T20:20:05.197915582Z",
    "isBlocked": true
}
```

### Remoção de documento

```shell
curl --location --request DELETE 'http://localhost:8080/api/document?id=106b8780-3bd7-4e73-a7df-726120ffc653'
```
