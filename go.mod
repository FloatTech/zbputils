module github.com/FloatTech/zbputils

go 1.19

require (
	github.com/FloatTech/floatbox v0.0.0-20230207075003-0f70b30c320d
	github.com/FloatTech/gg v1.1.2
	github.com/FloatTech/imgfactory v0.2.2-0.20230215052637-9f7b05520ca9
	github.com/FloatTech/rendercard v0.0.10-0.20230215092509-ff0745852f23
	github.com/FloatTech/sqlite v1.5.7
	github.com/FloatTech/ttl v0.0.0-20220715042055-15612be72f5b
	github.com/FloatTech/zbpctrl v1.5.3-0.20230130095145-714ad318cd52
	github.com/disintegration/imaging v1.6.2
	github.com/fumiama/cron v1.3.0
	github.com/fumiama/go-base16384 v1.6.4
	github.com/fumiama/go-registry v0.2.5
	github.com/go-playground/assert/v2 v2.2.0
	github.com/sirupsen/logrus v1.9.0
	github.com/tidwall/gjson v1.14.4
	github.com/wdvxdr1123/ZeroBot v1.6.8
)

require (
	github.com/RomiChan/syncx v0.0.0-20221202055724-5f842c53020e // indirect
	github.com/ericpauley/go-quantize v0.0.0-20200331213906-ae555eb2afa4 // indirect
	github.com/fumiama/go-simple-protobuf v0.1.0 // indirect
	github.com/fumiama/gofastTEA v0.0.10 // indirect
	github.com/fumiama/imgsz v0.0.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/image v0.3.0 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
	golang.org/x/text v0.6.0 // indirect
	modernc.org/libc v1.21.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.4.0 // indirect
	modernc.org/sqlite v1.20.0 // indirect
)

replace modernc.org/sqlite => github.com/fumiama/sqlite3 v1.20.0-with-win386

replace github.com/remyoudompheng/bigfft => github.com/fumiama/bigfft v0.0.0-20211011143303-6e0bfa3c836b
