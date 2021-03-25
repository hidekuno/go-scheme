module go-scheme/draw

go 1.16

replace local.packages/draw => ./

replace local.packages/fractal => ./fractal

require (
	github.com/hidekuno/go-scheme/scheme v0.0.0-20210325085929-bc906ad70210
	github.com/mattn/go-gtk v0.0.0-20191030024613-af2e013261f5
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	local.packages/draw v0.0.0-00010101000000-000000000000
	local.packages/fractal v0.0.0-00010101000000-000000000000
)
