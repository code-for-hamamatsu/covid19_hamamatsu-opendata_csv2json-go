# covid19_hamamatsu-opendata_csv2json-go

「[浜松市オープンデータカタログ](https://www.city.hamamatsu.shizuoka.jp/odpf/opendata/v1.html)」で公開されているCSVデータを  
「[浜松市新型コロナウイルス感染症対策サイト](https://stopcovid19.code4hamamatsu.org/)」で利用しているdata.jsonへ変換するプロジェクトです。  
(* Amazon API Gatewayから呼び出される AWS Lambda にデプロイして実行される事を想定した実装となっています。)

## サポートデータとID定義

| COVID-19サイトデータ名称 | COVID-19サイトデータID | 浜松市APIID |
| --- | --- | --- |
| 検査陽性者の状況 | main_summary | 221309_hamamatsu_covid19_patients & 221309_hamamatsu_covid19_patients_summary |
| 陽性患者の属性 | patients | 221309_hamamatsu_covid19_patients |
| 陽性患者数 | patients_summary | 221309_hamamatsu_covid19_patients |
| 検査実施人数 | inspection_persons| 221309_hamamatsu_covid19_test_people |
| 新型コロナウイルス感染症に関する相談件数| contacts | 221309_hamamatsu_covid19_call_center |

(* [検査陽性者の状況] は、死亡者のカウントのために、2つのCSVを参照しています。)

| 浜松市APIID | 浜松市オープンデータカタログページ |
| --- | --- |
| 221309_hamamatsu_covid19_patients | [コロナ 陽性患者の属性](https://www.city.hamamatsu.shizuoka.jp/odpf/opendata/v1.html?x=221309_hamamatsu_covid19_patients) |
| 221309_hamamatsu_covid19_patients_summary| [コロナ 陽性患者数](https://www.city.hamamatsu.shizuoka.jp/odpf/opendata/v1.html?x=221309_hamamatsu_covid19_patients_summary) |
| 221309_hamamatsu_covid19_test_people| [コロナ PCR検査実施状況](https://www.city.hamamatsu.shizuoka.jp/odpf/opendata/v1.html?x=221309_hamamatsu_covid19_test_people) |
| 221309_hamamatsu_covid19_call_center| [コロナ 相談件数](https://www.city.hamamatsu.shizuoka.jp/odpf/opendata/v1.html?x=221309_hamamatsu_covid19_call_center) |

## クエリパラメータ引数について

| key | value |
| --- | --- |
| type | GraphType-key:API-IDの配列 |

```bash
example
?type=main_summary:221309_hamamatsu_covid19_patients,main_summary:221309_hamamatsu_covid19_patients_summary,patients:221309_hamamatsu_covid19_patients,patients_summary:221309_hamamatsu_covid19_patients,inspection_persons:221309_hamamatsu_covid19_test_people,contacts:221309_hamamatsu_covid19_call_center
```

## Deploy to Lambda (zip)

```bash
commands
$ pwd
{workspaceRoot}/src/app

$ GOOS=linux GOARCH=amd64 go build -o ../../bin/main main.go
$ (cd ../../bin && zip -r ../lambda-package.zip *)
$ aws lambda update-function-code --function-name ${LAMBDA_FUNCTION_NAME} --zip-file fileb://../../lambda-package.zip
```

## Deploy to Lambda (Container image)

```bash
commands
$ pwd
{workspaceRoot}
$ aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com

$ docker build -f Dockerfile.release -t csv2json-release .
$ docker tag csv2json-release:latest ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/csv2json:latest
$ docker push ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/csv2json:latest
$ aws lambda update-function-code --function-name ${LAMBDA_FUNCTION_NAME} --image-uri ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/csv2json:latest
```

## Docker for local

```bash
commands
$ pwd
{workspaceRoot}

# build image & run container
$ docker build -f Dockerfile.debug -t csv2json-debug .
$ docker run --rm -p 9000:8080 csv2json-debug:latest /main

# request for test
$ curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}' -o ret.json

# request for test (with query parameters)
$ curl -XPOST  \
    "http://localhost:9000/2015-03-31/functions/function/invocations"  \
    -d '{ "queryStringParameters" : { "type" : "main_summary:5ab47071-3651-457c-ae2b-bfb8fdbe1af1" } }' \
    | jq -r .body
```

## ユニットテスト

```bash
commands
$ pwd
{workspaceRoot}/src/app

$ go test
```
