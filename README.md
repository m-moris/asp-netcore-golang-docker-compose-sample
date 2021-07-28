# ASP.NET Core  -> Golang 

ASP.NET Core WebAPI --> gRpc --> Golang のサンプル。

## golang スタンドアロン

スタンドアロンでの実行確認

#### 単体実行

```
cd go
go run main.go
grpcurl -plaintext -d '{"name" : "moris"}' localhost:50051 Greeter.SayHello
```

### docker での実行

```
docker build . -t morishima/aspgo-backend
docker run -p 50051:50051 -it morishima/aspgo-backend
```

### grpcurl

`grpcurl` で確認できる

```
grpcurl -plaintext -d '{"name" : "moris"}' localhost:50051 Greeter.SayHello
```

## asp-netcore

protoc からの自動生成は、Visual Studio 2019 が行い、protoファイルへはリンク参照だと後々エラーになるので、ファイルそのものを、プロジェクト配下にコピーしているので注意。

```
dotnet build
dotnet run
```

```
docker build . -t morishima/aspgo-forntend
docker run -p 8080:80 -it morishima/aspgo-frontend
```

VS上でデバッグして確認してもよい。

## docker compose 

2つ起動して、asp.net core 側から gRpc サーバーを呼び出す。

```json
version: '3.4'

services:
  backend:
    image: morishima/aspgo-backend
    build:
      context: ./go
      dockerfile: Dockerfile
    environment:
      - GRPC_PORT=50051
    ports:
      - "50051:50051"
  frontend:
    image: morishima/aspgo-frontend
    build:
      context: ./asp-netcore
      dockerfile: Dockerfile
    environment:
      - GRPCOPTIONS__SERVER=backend
      - GRPCOPTIONS__PORT=50051
      - ASPNETCORE_URLS=http://+:80
    ports:
      - "8080:80"
    volumes:
      - ~/.aspnet/https:/https:ro
```

以下でビルド、実行。


```
docker-compose build
docker-compose up
docker-compose down
curl http://localhost:8080/hello
```

## Azure WebApps for **Multiple** Container 


Docker hub にログインしつつ、タグ付けし、プッシュする。

```
docker login
docker-compose -f azure-compose.yml push
```

初期設定とリソースグループの作成

```
rg=webapps-test
az configure --defaults location=japaneast
az configure --defaults group=$rg
az group create -n $rg
```

App Service Plan の作成

```
plan=containerplan
az appservice plan create -n $plan --sku S1 --is-linux
```

App Service の作成とデプロイ。既存のリソースがある場合は、設定を上書きしてくれる。

```
az webapp create --plan $plan --name mycontainer20210808 \
                 --multicontainer-config-type compose \
                 --multicontainer-config-file docker-compose.yml
```

しばらく待つ。Kudu で LogStream を見ると状況が分る。

APIを叩く。

```
curl http://mycontainer20210808.azurewebsites.net/hello
```