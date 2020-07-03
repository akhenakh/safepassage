module github.com/akhenakh/safepassage

go 1.14

require (
	github.com/ProtonMail/gopenpgp/v2 v2.0.1
	github.com/namsral/flag v1.7.4-pre
)

replace golang.org/x/crypto => github.com/ProtonMail/crypto v0.0.0-20200416114516-1fa7f403fb9c
