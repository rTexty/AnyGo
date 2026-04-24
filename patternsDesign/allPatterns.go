package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

type Resource string
type Item string

const (
	Wood  Resource = "WOOD"
	Herb  Resource = "HERB"
	Metal Resource = "METAL"
)

const (
	Bow   Item = "BOW"
	Gun   Item = "GUN"
	Armor Item = "ARMOR"
	None  Item = "NONE"
)

type CreatureTypeData struct {
	Name                          string
	MaxHP, MaxStamina, BaseDamage int
	HarvestStrategy               Harvester
}

var mobs = map[string]CreatureTypeData{
	"Navi": {
		Name:            "Navi",
		MaxHP:           120,
		MaxStamina:      100,
		BaseDamage:      12,
		HarvestStrategy: NaviHarvester{},
	},
	"Human": {
		Name:            "Human",
		MaxHP:           100,
		MaxStamina:      120,
		BaseDamage:      10,
		HarvestStrategy: HumanHarvester{},
	},
	"Ikran": {
		Name:            "Ikran",
		MaxHP:           80,
		MaxStamina:      150,
		BaseDamage:      8,
		HarvestStrategy: ForbiddenHarvester{},
	},
	"Direhorse": {
		Name:            "Direhorse",
		MaxHP:           90,
		MaxStamina:      130,
		BaseDamage:      6,
		HarvestStrategy: ForbiddenHarvester{},
	},
	"Thanator": {
		Name:            "Thanator",
		MaxHP:           150,
		MaxStamina:      110,
		BaseDamage:      18,
		HarvestStrategy: ForbiddenHarvester{},
	},
}
var once sync.Once
var instance *Game

// Stock is a structural Adapter interface for working with resource storage.
type Stock interface {
	Get(Resource) int
	Add(Resource, int)
	Take(Resource, int) bool
}

type StockAdapter struct {
	data map[Resource]int
}

func NewStockAdapter(wood, herb, metal int) *StockAdapter {
	return &StockAdapter{
		data: map[Resource]int{
			Wood:  wood,
			Herb:  herb,
			Metal: metal,
		},
	}
}

func (s *StockAdapter) Get(r Resource) int {
	return s.data[r]
}

func (s *StockAdapter) Add(r Resource, amount int) {
	s.data[r] += amount
}

func (s *StockAdapter) Take(r Resource, amount int) bool {
	if s.data[r] < amount {
		return false
	}
	s.data[r] -= amount
	return true
}

type Zone struct {
	name      string
	stock     Stock
	creatures map[string]bool
}
type Creature struct {
	id          string
	type_       CreatureTypeData
	currHp      int
	currStamina int
	resources   map[Resource]int
	items       map[Item]int
	weapon      Item
	armor       Item
	partner     string
	limits      *CreatureTypeData
	zoneName    string
}

type Game struct {
	zones     map[string]*Zone
	creatures map[string]*Creature
}

func GetInstance() *Game {
	once.Do(func() {
		instance = &Game{
			zones:     make(map[string]*Zone),
			creatures: make(map[string]*Creature),
		}
	})
	return instance
}

func (g *Game) CreateZone(name string, wood, herb, metal int) string {
	g.zones[name] = &Zone{
		name:      name,
		stock:     NewStockAdapter(wood, herb, metal),
		creatures: make(map[string]bool),
	}
	return "OK"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func newCreature(id, typeName, zoneName string, hp, stamina, wood, herb, metal int) *Creature {
	data := mobs[typeName]
	limits := data
	return &Creature{
		id:          id,
		type_:       data,
		currHp:      min(hp, data.MaxHP),
		currStamina: min(stamina, data.MaxStamina),
		resources: map[Resource]int{
			Wood:  wood,
			Herb:  herb,
			Metal: metal,
		},
		items: map[Item]int{
			Bow:   0,
			Gun:   0,
			Armor: 0,
		},
		weapon:   None,
		armor:    None,
		partner:  "",
		limits:   &limits,
		zoneName: zoneName,
	}
}

func (g *Game) CreateCreature(id, typeName, zoneName string, hp, stamina, wood, herb, metal int) string {
	c := newCreature(id, typeName, zoneName, hp, stamina, wood, herb, metal)
	g.creatures[id] = c
	g.zones[zoneName].creatures[id] = true
	return "OK"
}

// Clone is the Prototype pattern implementation.
func (c *Creature) Clone(newID string) *Creature {
	clone := *c
	clone.id = newID
	clone.partner = ""
	clone.resources = map[Resource]int{
		Wood:  c.resources[Wood],
		Herb:  c.resources[Herb],
		Metal: c.resources[Metal],
	}
	clone.items = map[Item]int{
		Bow:   c.items[Bow],
		Gun:   c.items[Gun],
		Armor: c.items[Armor],
	}
	return &clone
}

func (g *Game) CloneCreature(newID, existingID string) string {
	src := g.creatures[existingID]
	c := src.Clone(newID)
	g.creatures[newID] = c
	g.zones[c.zoneName].creatures[newID] = true
	return "OK"
}

func (g *Game) Give(id string, item Item, count int) string {
	c := g.creatures[id]
	if c.currHp <= 0 {
		return "ERR DEAD"
	}
	c.items[item] += count
	return "OK"
}

func (g *Game) Move(id string, zoneName string) string {
	c := g.creatures[id]
	if c.currHp <= 0 {
		return "ERR DEAD"
	}

	if c.partner != "" {
		p := g.creatures[c.partner]
		if p.currHp <= 0 {
			return "ERR DEAD"
		}
		if c.currStamina < 1 || p.currStamina < 1 {
			return "ERR RESOURCES"
		}
		oldCZone := g.zones[c.zoneName]
		oldPZone := g.zones[p.zoneName]
		delete(oldCZone.creatures, c.id)
		delete(oldPZone.creatures, p.id)
		c.zoneName = zoneName
		p.zoneName = zoneName
		g.zones[zoneName].creatures[c.id] = true
		g.zones[zoneName].creatures[p.id] = true
		c.currStamina--
		p.currStamina--
		return "OK"
	}

	if c.currStamina < 1 {
		return "ERR RESOURCES"
	}
	delete(g.zones[c.zoneName].creatures, c.id)
	c.zoneName = zoneName
	g.zones[zoneName].creatures[c.id] = true
	c.currStamina--
	return "OK"
}

func (g *Game) Harvest(id string, resource Resource, amount int) string {
	c := g.creatures[id]
	if c.currHp <= 0 {
		return "ERR DEAD"
	}
	if !c.type_.HarvestStrategy.CanHarvest(resource) {
		return "ERR FORBIDDEN"
	}
	z := g.zones[c.zoneName]
	if c.currStamina < amount || z.stock.Get(resource) < amount {
		return "ERR RESOURCES"
	}
	z.stock.Take(resource, amount)
	c.resources[resource] += amount
	c.currStamina -= amount
	return "OK"
}

func (g *Game) Equip(id string, item Item) string {
	c := g.creatures[id]
	if c.currHp <= 0 {
		return "ERR DEAD"
	}
	switch item {
	case Bow:
		if c.type_.Name != "Navi" {
			return "ERR FORBIDDEN"
		}
	case Gun:
		if c.type_.Name != "Human" {
			return "ERR FORBIDDEN"
		}
	case Armor:
		if c.type_.Name != "Navi" && c.type_.Name != "Human" {
			return "ERR FORBIDDEN"
		}
	default:
		return "ERR FORBIDDEN"
	}
	if c.items[item] < 1 {
		return "ERR RESOURCES"
	}
	if item == Armor {
		c.armor = Armor
	} else {
		c.weapon = item
	}
	return "OK"
}

func weaponBonus(i Item) int {
	switch i {
	case Bow:
		return 5
	case Gun:
		return 8
	default:
		return 0
	}
}

func (g *Game) clearLink(c *Creature) {
	if c.partner == "" {
		return
	}
	p := g.creatures[c.partner]
	c.partner = ""
	if p != nil {
		p.partner = ""
	}
}

func (g *Game) Attack(attackerID, defenderID string) string {
	attacker := g.creatures[attackerID]
	defender := g.creatures[defenderID]

	if attacker.currHp <= 0 || defender.currHp <= 0 {
		return "ERR DEAD"
	}
	if attacker.zoneName != defender.zoneName {
		return "ERR LOCATION"
	}
	if attacker.currStamina < 5 {
		return "ERR RESOURCES"
	}

	damage := attacker.type_.BaseDamage + weaponBonus(attacker.weapon)
	if defender.armor == Armor {
		damage -= 3
	}
	if damage < 0 {
		damage = 0
	}

	attacker.currStamina -= 5
	defender.currHp -= damage
	if defender.currHp < 0 {
		defender.currHp = 0
	}

	if defender.currHp == 0 {
		g.clearLink(defender)
		z := g.zones[defender.zoneName]
		z.stock.Add(Wood, defender.resources[Wood])
		z.stock.Add(Herb, defender.resources[Herb])
		z.stock.Add(Metal, defender.resources[Metal])
		defender.resources[Wood] = 0
		defender.resources[Herb] = 0
		defender.resources[Metal] = 0
	}

	return "OK"
}

func (g *Game) Link(naviID, beastID string) string {
	navi := g.creatures[naviID]
	beast := g.creatures[beastID]

	if navi.currHp <= 0 || beast.currHp <= 0 {
		return "ERR DEAD"
	}
	if navi.zoneName != beast.zoneName {
		return "ERR LOCATION"
	}
	if navi.type_.Name != "Navi" {
		return "ERR FORBIDDEN"
	}
	if beast.type_.Name != "Ikran" && beast.type_.Name != "Direhorse" {
		return "ERR FORBIDDEN"
	}
	if navi.partner != "" || beast.partner != "" {
		return "ERR FORBIDDEN"
	}
	navi.partner = beast.id
	beast.partner = navi.id
	return "OK"
}

func (g *Game) Unlink(id string) string {
	c := g.creatures[id]
	if c.currHp <= 0 {
		return "ERR DEAD"
	}
	if c.partner == "" {
		return "ERR FORBIDDEN"
	}
	g.clearLink(c)
	return "OK"
}

func (g *Game) Status(id string) string {
	c := g.creatures[id]
	partner := c.partner
	if partner == "" {
		partner = "NONE"
	}
	weapon := string(c.weapon)
	armor := string(c.armor)
	return fmt.Sprintf(
		"OK %s %s %s HP %d/%d ST %d/%d LINK %s EQUIP %s %s",
		c.id, c.type_.Name, c.zoneName, c.currHp, c.limits.MaxHP, c.currStamina, c.limits.MaxStamina, partner, weapon, armor,
	)
}

func (g *Game) Inv(id string) string {
	c := g.creatures[id]
	return fmt.Sprintf(
		"OK %s INV %d %d %d ITEMS %d %d %d",
		c.id,
		c.resources[Wood], c.resources[Herb], c.resources[Metal],
		c.items[Bow], c.items[Gun], c.items[Armor],
	)
}

func (g *Game) ZoneStatus(zoneName string) string {
	z := g.zones[zoneName]
	ids := make([]string, 0, len(z.creatures))
	for id := range z.creatures {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	if len(ids) == 0 {
		return fmt.Sprintf(
			"OK %s RES %d %d %d CREATURES 0",
			zoneName, z.stock.Get(Wood), z.stock.Get(Herb), z.stock.Get(Metal),
		)
	}
	return fmt.Sprintf(
		"OK %s RES %d %d %d CREATURES %d %s",
		zoneName, z.stock.Get(Wood), z.stock.Get(Herb), z.stock.Get(Metal), len(ids), strings.Join(ids, " "),
	)
}

func parseResource(s string) Resource {
	switch s {
	case "Wood":
		return Wood
	case "Herb":
		return Herb
	default:
		return Metal
	}
}

func parseItem(s string) Item {
	switch s {
	case "BOW":
		return Bow
	case "GUN":
		return Gun
	default:
		return Armor
	}
}

type Harvester interface {
	CanHarvest(r Resource) bool
}
type NaviHarvester struct{}
type HumanHarvester struct{}
type ForbiddenHarvester struct{}

func (h ForbiddenHarvester) CanHarvest(r Resource) bool {
	return false
}

func (h HumanHarvester) CanHarvest(r Resource) bool {
	if r == Metal {
		return true
	}
	return false
}
func (h NaviHarvester) CanHarvest(r Resource) bool {
	if r == Wood || r == Herb {
		return true
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	game := GetInstance()
	var q int
	fmt.Fscan(in, &q)

	for i := 0; i < q; i++ {
		var cmd string
		fmt.Fscan(in, &cmd)
		switch cmd {
		case "ZONE":
			var zoneName string
			var wood, herb, metal int
			fmt.Fscan(in, &zoneName, &wood, &herb, &metal)
			fmt.Fprintln(out, game.CreateZone(zoneName, wood, herb, metal))
		case "CREATURE":
			var id, typeName, zoneName string
			var hp, stamina, wood, herb, metal int
			fmt.Fscan(in, &id, &typeName, &zoneName, &hp, &stamina, &wood, &herb, &metal)
			fmt.Fprintln(out, game.CreateCreature(id, typeName, zoneName, hp, stamina, wood, herb, metal))
		case "CLONE":
			var newID, existingID string
			fmt.Fscan(in, &newID, &existingID)
			fmt.Fprintln(out, game.CloneCreature(newID, existingID))
		case "GIVE":
			var id, itemName string
			var count int
			fmt.Fscan(in, &id, &itemName, &count)
			fmt.Fprintln(out, game.Give(id, parseItem(itemName), count))
		case "MOVE":
			var id, zoneName string
			fmt.Fscan(in, &id, &zoneName)
			fmt.Fprintln(out, game.Move(id, zoneName))
		case "HARVEST":
			var id, resourceName string
			var amount int
			fmt.Fscan(in, &id, &resourceName, &amount)
			fmt.Fprintln(out, game.Harvest(id, parseResource(resourceName), amount))
		case "EQUIP":
			var id, itemName string
			fmt.Fscan(in, &id, &itemName)
			fmt.Fprintln(out, game.Equip(id, parseItem(itemName)))
		case "ATTACK":
			var attackerID, defenderID string
			fmt.Fscan(in, &attackerID, &defenderID)
			fmt.Fprintln(out, game.Attack(attackerID, defenderID))
		case "LINK":
			var naviID, beastID string
			fmt.Fscan(in, &naviID, &beastID)
			fmt.Fprintln(out, game.Link(naviID, beastID))
		case "UNLINK":
			var id string
			fmt.Fscan(in, &id)
			fmt.Fprintln(out, game.Unlink(id))
		case "STATUS":
			var id string
			fmt.Fscan(in, &id)
			fmt.Fprintln(out, game.Status(id))
		case "INV":
			var id string
			fmt.Fscan(in, &id)
			fmt.Fprintln(out, game.Inv(id))
		case "ZONESTATUS":
			var zoneName string
			fmt.Fscan(in, &zoneName)
			fmt.Fprintln(out, game.ZoneStatus(zoneName))
		}
	}
}
