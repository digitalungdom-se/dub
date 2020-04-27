package assets

type contributor struct {
	Name      string
	Github    string
	username  string
	DiscordID string
	Role      ContributorType
}

type ContributorType int8

const (
	Maintainer ContributorType = iota
	Contributor
)

var kelszo = contributor{Name: "Kelvin Szolnoky", Github: "https://github.com/kelszo", DiscordID: "217632464531619852", Role: Maintainer}
var simonthor = contributor{Name: "Simon Thor", Github: "https://github.com/simonthor", DiscordID: "228889878861971456", Role: Contributor}
var zigolox = contributor{Name: "Simon Sondén", Github: "https://github.com/Zigolox", DiscordID: "384331517243031552", Role: Contributor}

var Contributors = []contributor{kelszo, simonthor, zigolox}
