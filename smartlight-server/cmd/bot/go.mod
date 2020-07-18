module bot

go 1.14

require (
	github.com/go-co-op/gocron v0.2.0
	github.com/julienschmidt/httprouter v1.3.0
	gobot.io/x/gobot v1.14.0
	internal/schedule v0.0.0
	internal/usage v0.0.0
)

replace internal/schedule => ./../../internal/schedule

replace internal/usage => ./../../internal/usage
