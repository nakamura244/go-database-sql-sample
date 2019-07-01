# このリポジトリに関して
実際のmysqlを使わずにunit testを実現するためのサンプルコード

see detail : https://tsuyoshi-nakamura.hatenablog.com/entry/2019/07/01/094425


# ディレクトリ構成に関して
```
├── db                              ... db操作のcore部分
│   ├── iface
│   │   └── sql.go                  ... 利用するpkg/database/sql のメソッドのinterface登録
│   ├── sql.go                      ... db操作のメソッド定義
│   └── sql_test.go
├── interfaces                      ... db.sqlとrepositoryディレクトリとつなぎ役
│   ├── sql_handler.go              ... db.sqlで定義したメソッドのinterface登録
│   ├── sql_repository.go           ... db.sqlで定義したメソッドのinterfaceを使ってメソッドを定義
│   └── sql_repository_test.go
├── main.go
└── repository                      ... interfaces.Ssql_repository.goで定義したメソッドのinterface登録
    └── db_repository.go

```

# Checked version of Go
1.12.4


# How to run
1. clone
1. `go test -v -cover ./...`