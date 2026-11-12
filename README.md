# e2e

複数のバックエンドサービスをまたいで動作を確認する E2E テストのリポジトリです。

## テストの種類と線引き

| 種類 | 置き場所 | 書くもの | 書かないもの |
|---|---|---|---|
| endpoint | `tests/endpoint/<service>/` | 実際の環境で、サービス間の連携が起きる単一の操作(例: 申込を承認すると口座が作られる) | 入力値チェックなど、1つのサービスの中で完結する分岐(各サービスの UT/IT で確認する) |
| journey | `tests/journey/` | 重要な業務フローを通しで確認するシナリオ | 細かい分岐の網羅 |

迷ったら「各サービスの IT で確認できるか」を考えてください。確認できるなら、E2E には書きません。

## データ方針

- テストデータは、すべて `support/fixture` から**各サービスの公開 API** を通して作成します。DB への直接投入はしません。
- 値には `support/idgen` を使い、実行ID(`E2E_RUN_ID`)を含む一意な値にします。これで並列実行や他の実行と衝突しません。
- テストデータは削除しません。実行IDで分離されています。
- **マスタデータは、実行先の環境に投入済みであることが前提です**(各サービスのマイグレーションなどで管理)。E2E 側では作成も削除もしません。
- **外部サービスはスタブしません。** 外部サービスとの連携を含むフローは、外部が本物かサンドボックスにつながっている環境でのみ確認します。

## 非同期処理の待ち方

キューやイベント経由の処理結果は、`assertion.Eventually` で待ちます。`time.Sleep` は使いません。
待ち時間の上限と再試行の間隔は、`envs/<env>.yaml` の `eventually` で環境ごとに設定します。

## ディレクトリ構成

```
support/
  client/     各サービスの API クライアント(現在は仮の手書き実装。OpenAPI からの生成に置き換える予定)
  fixture/    API 経由でテストデータを作るヘルパー
  assertion/  Eventually など
  idgen/      一意な値と実行ID
  config/     envs/*.yaml の読み込み
tests/
  endpoint/   サービス間の連携が起きる単一の操作
  journey/    業務フロー
envs/         環境ごとの接続先設定
```

## 実行方法

テストには build tag `e2e` を付けています。タグなしの `go test ./...` では実行されません。

```sh
task test                    # local 環境で全テスト
task test E2E_ENV=stg        # stg 環境で全テスト
task test:endpoint
task test:journey
task test -- -run TestJourney_AccountOpening -v
```

Task を使わない場合は次のとおりです。

```sh
E2E_ENV=stg go test -tags e2e -count=1 ./tests/...
```

### 環境変数

| 変数 | 説明 | 既定値 |
|---|---|---|
| `E2E_ENV` | 読み込む `envs/<env>.yaml` | `local` |
| `E2E_RUN_ID` | 実行ID。テストデータと `X-E2E-Run-Id` ヘッダに付与される | 実行時刻+乱数 |
| `E2E_<SERVICE>_BASE_URL` | 接続先の上書き(例: `E2E_ACCOUNT_BASE_URL`) | - |
| `E2E_CONFIG_DIR` | `envs/` の場所の上書き | リポジトリの `envs/` |

### エディタ設定

gopls が build tag 付きのファイルを解析できるよう、`.vscode/settings.json` で `-tags=e2e` を指定しています。

## CI

`.gitlab-ci.yml` の `e2e` ジョブは、次のときに実行されます。

- 各サービスのリポジトリから trigger で起動されたとき
- 定期実行
- 手動実行
- このリポジトリのデフォルトブランチが更新されたとき

各サービスのリポジトリからは、次のように起動します。

```yaml
e2e:
  stage: e2e
  trigger:
    project: <group>/e2e
    branch: main
    strategy: depend   # E2E の結果を呼び出し元のパイプラインに反映する
```
