package evediscordkm

type idNameEntity struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Attacker struct {
	Character   idNameEntity `json:"character"`
	Corporation idNameEntity `json:"corporation"`
	Alliance    idNameEntity `json:"alliance"`
	DamageDone  int          `json:"damageDone"`
	FinalBlow   bool         `json:"finalBlow"`
	ShipType    idNameEntity `json:"shipType"`
	WeaponType  idNameEntity `json:"weaponType"`
}

type Victim struct {
	Character   idNameEntity `json:"character"`
	Corporation idNameEntity `json:"corporation"`
	Alliance    idNameEntity `json:"alliance"`
	DamageTaken int          `json:"damageTalen"`
	ShipType    idNameEntity `json:"shipType"`
}

type KillMail struct {
	AttackerCount int          `json:"attackerCount"`
	Attackers     []Attacker   `json:"attackers"`
	Victim        Victim       `json:"victim"`
	SolarSystem   idNameEntity `json:"solarSystem"`
}

type Metadata struct {
	TotalValue float64 `json:"totalValue"`
}

type KillReport struct {
	Killmail KillMail `json:"killmail"`
	KillID   int      `json:"killID"`
	Metadata Metadata `json:"zkb"`
}

type byDamage []Attacker

func (arr byDamage) Len() int      { return len(arr) }
func (arr byDamage) Swap(x, y int) { arr[x], arr[y] = arr[y], arr[x] }
func (arr byDamage) Less(x, y int) bool {
	return arr[x].FinalBlow || (arr[x].DamageDone < arr[y].DamageDone)
}
