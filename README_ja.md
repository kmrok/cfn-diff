# cfn-diff

`cfn-diff` はAWS CloudFormationリソースに対するdiffツールです。

## インストール

**go get**

```sh
$ go get github.com/kmrok/cfn-diff/cmd/cfn-diff
```

**build**

```sh
$ git clone https://github.com/kmrok/cfn-diff.git
$ make install
```

## はじめに

AWSのクレデンシャルを環境変数に設定してください。

```sh
$ export AWS_ACCESS_KEY_ID=XXXXXX
$ export AWS_SECRET_ACCESS_KEY=XXXXXX
$ export AWS_REGION=XXXXXX
```

設定ファイル `.cfndiff.yml` ( もしくは `cfndiff.yml` ) をリポジトリに追加してください。

```yaml
# .cfndiff.yml

run:
  # 差分検出時に異常終了する場合、trueに設定 (任意)
  enable_ci_mode: true

# ドリフト検出したいスタックの名前を記述 (任意)
# 正規表現使用可能
# デフォルト: 全てのスタック
stack_with_drift_detection:
  - stg-.*

# ドリフト検出したくないスタックの名前を記述 (任意)
# 正規表現使用可能
stack_without_drift_detection:
  - stg-aurora
  - stg-.*-ecs-task-def
  - stg-alb

# スタック名とテンプレートファイルの対応を記述 (任意)
stack_template_maps:
  - stack_name: stg-ecr
    template_name: ./ecr.yml
  - stack_name: stg-ecs-task-def
    template_name: ./ecs-task-def.yml
```

## 使用方法

### ドリフト検出

```sh
$ cfn-diff drift
$ cfn-diff drift -f config.yml
```

### 変更セット差分検出

```sh
$ cfn-diff changeset
$ cfn-diff changeset -f config.yml
```
