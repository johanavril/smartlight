module server

go 1.14

require (
	github.com/go-co-op/gocron v0.2.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.7.0
	internal/connection v0.0.0
	internal/lamp v0.0.0
	internal/schedule v0.0.0
	internal/setting v0.0.0
	internal/usage v0.0.0
)

replace internal/connection => ./../../internal/connection

replace internal/lamp => ./../../internal/lamp

replace internal/schedule => ./../../internal/schedule

replace internal/setting => ./../../internal/setting

replace internal/usage => ./../../internal/usage