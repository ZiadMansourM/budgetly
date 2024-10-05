## ToDo
- [ ] Finish Authentication and Authorization. Plus, OAuth.
- [ ] Add Templates Support. Explore htmgo.
- [ ] Support TXs for Database Operations.
- [ ] Admin Dashboard.

```bash
ziadh@Ziads-MacBook-Air budgetly % tree
.
├── README.md
├── budgetly
│   ├── cmd
│   │   ├── api
│   │   │   └── server_builder.go
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   └── apps
│   │       └── users
│   │           ├── app.go
│   │           ├── handlers.go
│   │           ├── models.go
│   │           └── services.go
│   ├── pkg
│   │   ├── db
│   │   │   └── connection.go
│   │   └── middlewares
│   │       └── logging.go
│   └── utils
│       └── utils.go
├── compose.yaml
└── scripts
    └── initdb.sql

11 directories, 14 files
```
