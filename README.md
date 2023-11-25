# mocha

![mocha](mocha.png)
#### meow-meow REST server

## export & proto build

```bash
export PATH="/Users/{user_go_path}/go/bin:$PATH"
protoc -I . --go_out=./proto --go-grpc_out=./proto ./proto/*.proto
```

## run
```bash
go run ./cmd/server
```

## view message table list
```bash
export AWS_ACCESS_KEY_ID=x
export AWS_SECRET_ACCESS_KEY=x
aws dynamodb list-tables --endpoint-url http://localhost:8008
```
### query messages
```bash
aws dynamodb query \
    --table-name messages \
    --key-condition-expression "channel_id = :val" \
    --expression-attribute-values '{":val": {"N": "1"}}' \
    --limit 1 \
    --endpoint-url http://localhost:8008
```
### query last
```bash
aws dynamodb query \
    --table-name messages \
    --key-condition-expression "channel_id = :val" \
    --expression-attribute-values '{":val": {"N": "1"}}' \
    --exclusive-start-key "{ \"channel_id\": {\"N\": \"1\"}, \"id\": {\"N\": \"message_id\"} }" \
    --endpoint-url http://localhost:8008
```