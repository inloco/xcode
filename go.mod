module github.com/inloco/xcode

go 1.17

require (
	github.com/github/smimesign v0.2.0
)

replace (
	github.com/github/smimesign v0.2.0 => github.com/inloco/smimesign v0.2.1-0.20220322075834-64a83fea8d33
)
