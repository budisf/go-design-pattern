# Tolong dibaca ya!

Ini merupakan sedikit *introduction* notes untuk project ethical untuk menggunakan tools [golang-migrate][1].

Kalau kalian tertarik untuk buat readme documentation seperti ini, bisa cek disini: [Markdown Cheatsheet](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet)

[1]: https://github.com/golang-migrate/migrate

## Concern
untuk sekarang masih belum bisa import as library disini\
karena belum paham penggunaan driver secara build in function\
[Use in our golang project](https://github.com/golang-migrate/migrate#use-in-your-go-project)\
kalau ada yang sudah bisa silahkan di implementasi
 
## Tools yang diperlukan
- [Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)\
pengecekan version setelah download bisa melakukan\
`migrate version`

## Migrate Script
Berikut ini adalah script-script contoh yang umum digunakan di cli
- `migrate create -ext sql -dir app/databases/migrations -seq init_schema`
- `migrate -path app/databases/migrations -database "postgres://root:secret@localhost:5433/econolab_ethical?sslmode=disable" up`
- `migrate -path app/databases/migrations -database "postgres://root:secret@localhost:5433/econolab_ethical?sslmode=disable" down`

## Cara Penggunaan
Migrate memiliki 2 file ketika di generate, yaitu file up dan file down,\
cara untuk naming migrations bisa dilihat disini\
[Best practices: How to write migrations](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md)\
Untuk kenapa migration punya up and down bisa cek [disini](https://github.com/golang-migrate/migrate/blob/master/FAQ.md#why-two-separate-files-up-and-down-for-a-migration)


## Makefile
Untuk penggunaan migrate ini kita kordinasikan dengan penggunaan docker untuk mempermudah proses *build up* database, scriptnya bisa dilihat dalam Makefile mengikuti [tutorial](https://www.youtube.com/watch?v=0CYkrGIJkpw) ini.

command yang tersedia di makefile berupa:
- make postgres-start
- make createdb
- make dropdb
- make postgres-stop
- make postgres-delete
- make migrate-up
- make migrate-down