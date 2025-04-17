module go-scheme/web

go 1.23.0

replace local.packages/web => ./

require (
	github.com/gorilla/sessions v1.2.1
	github.com/hidekuno/go-scheme/scheme v0.0.0-20210325085929-bc906ad70210
	golang.org/x/net v0.38.0
	local.packages/web v0.0.0-00010101000000-000000000000
)

require github.com/gorilla/securecookie v1.1.1 // indirect
