// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package item

import (
	"encoding/json"
	"fmt"
)

type Item_e int

const (
	None Item_e = iota
	Adze
	Arbalest
	Arrows
	Axes
	Backpack
	Ballistae
	Bark
	Barrel
	Bladder
	Blubber
	Boat
	BoneArmour
	Bones
	Bows
	Bread
	Breastplate
	Candle
	Canoes
	Carpets
	Catapult
	Cattle
	Cauldrons
	Chain
	China
	Clay
	Cloth
	Clubs
	Coal
	Coffee
	Coins
	Cotton
	Cuirass
	Cuirboilli
	Diamond
	Diamonds
	Drum
	Elephant
	Falchion
	Fish
	Flax
	Flour
	Flute
	Fodder
	Frame
	Frankincense
	Fur
	Glasspipe
	Goats
	Gold
	Grain
	Grape
	Gut
	HBow
	Harp
	Haube
	Heaters
	Helm
	Herbs
	Hive
	Hoe
	Honey
	Hood
	Horn
	Horses
	Jade
	Jerkin
	Kayak
	Ladder
	Leather
	Logs
	Lute
	Mace
	Mattock
	Metal
	MillStone
	Musk
	Net
	Oar
	Oil
	Olives
	Opium
	Ores
	Paddle
	Palanquin
	Parchment
	Pavis
	Pearls
	Pellets
	People
	Pewter
	Picks
	Plows
	Provisions
	Quarrel
	Rake
	Ram
	Ramp
	Ring
	Rope
	Rug
	Saddle
	Saddlebag
	Salt
	Sand
	Scale
	Sculpture
	Scutum
	Scythe
	Shackle
	Shaft
	Shield
	Shovel
	Silk
	Silver
	Skin
	Slaves
	Slings
	Snare
	Spear
	Spetum
	Spice
	Statue
	Stave
	Stones
	String
	Sugar
	Sword
	Tapestries
	Tea
	Tobacco
	Trap
	Trews
	Trinket
	Trumpet
	Urn
	Wagons
	Wax
)

// MarshalJSON implements the json.Marshaler interface.
func (e Item_e) MarshalJSON() ([]byte, error) {
	return json.Marshal(EnumToString[e])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *Item_e) UnmarshalJSON(data []byte) error {
	var s string
	var ok bool
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	} else if *e, ok = StringToEnum[s]; !ok {
		return fmt.Errorf("invalid Item %q", s)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e Item_e) String() string {
	if str, ok := EnumToString[e]; ok {
		return str
	}
	return fmt.Sprintf("Item(%d)", int(e))
}

var (
	// EnumToString is a helper map for marshalling the enum
	EnumToString = map[Item_e]string{
		None:         "N/A",
		Adze:         "Adze",
		Arbalest:     "Arbalest",
		Arrows:       "Arrows",
		Axes:         "Axes",
		Backpack:     "Backpack",
		Ballistae:    "Ballistae",
		Bark:         "Bark",
		Barrel:       "Barrel",
		Bladder:      "Bladder",
		Blubber:      "Blubber",
		Boat:         "Boat",
		BoneArmour:   "BoneArmour",
		Bones:        "Bones",
		Bows:         "Bows",
		Bread:        "Bread",
		Breastplate:  "Breastplate",
		Candle:       "Candle",
		Canoes:       "Canoes",
		Carpets:      "Carpets",
		Catapult:     "Catapult",
		Cattle:       "Cattle",
		Cauldrons:    "Cauldrons",
		Chain:        "Chain",
		China:        "China",
		Clay:         "Clay",
		Cloth:        "Cloth",
		Clubs:        "Clubs",
		Coal:         "Coal",
		Coffee:       "Coffee",
		Coins:        "Coins",
		Cotton:       "Cotton",
		Cuirass:      "Cuirass",
		Cuirboilli:   "Cuirboilli",
		Diamond:      "Diamond",
		Diamonds:     "Diamonds",
		Drum:         "Drum",
		Elephant:     "Elephant",
		Falchion:     "Falchion",
		Fish:         "Fish",
		Flax:         "Flax",
		Flour:        "Flour",
		Flute:        "Flute",
		Fodder:       "Fodder",
		Frame:        "Frame",
		Frankincense: "Frankincense",
		Fur:          "Fur",
		Glasspipe:    "Glasspipe",
		Goats:        "Goats",
		Gold:         "Gold",
		Grain:        "Grain",
		Grape:        "Grape",
		Gut:          "Gut",
		HBow:         "HBow",
		Harp:         "Harp",
		Haube:        "Haube",
		Heaters:      "Heaters",
		Helm:         "Helm",
		Herbs:        "Herbs",
		Hive:         "Hive",
		Hoe:          "Hoe",
		Honey:        "Honey",
		Hood:         "Hood",
		Horn:         "Horn",
		Horses:       "Horses",
		Jade:         "Jade",
		Jerkin:       "Jerkin",
		Kayak:        "Kayak",
		Ladder:       "Ladder",
		Leather:      "Leather",
		Logs:         "Logs",
		Lute:         "Lute",
		Mace:         "Mace",
		Mattock:      "Mattock",
		Metal:        "Metal",
		MillStone:    "MillStone",
		Musk:         "Musk",
		Net:          "Net",
		Oar:          "Oar",
		Oil:          "Oil",
		Olives:       "Olives",
		Opium:        "Opium",
		Ores:         "Ores",
		Paddle:       "Paddle",
		Palanquin:    "Palanquin",
		Parchment:    "Parchment",
		Pavis:        "Pavis",
		Pearls:       "Pearls",
		Pellets:      "Pellets",
		People:       "People",
		Pewter:       "Pewter",
		Picks:        "Picks",
		Plows:        "Plows",
		Provisions:   "Provisions",
		Quarrel:      "Quarrel",
		Rake:         "Rake",
		Ram:          "Ram",
		Ramp:         "Ramp",
		Ring:         "Ring",
		Rope:         "Rope",
		Rug:          "Rug",
		Saddle:       "Saddle",
		Saddlebag:    "Saddlebag",
		Salt:         "Salt",
		Sand:         "Sand",
		Scale:        "Scale",
		Sculpture:    "Sculpture",
		Scutum:       "Scutum",
		Scythe:       "Scythe",
		Shackle:      "Shackle",
		Shaft:        "Shaft",
		Shield:       "Shield",
		Shovel:       "Shovel",
		Silk:         "Silk",
		Silver:       "Silver",
		Skin:         "Skin",
		Slaves:       "Slaves",
		Slings:       "Slings",
		Snare:        "Snare",
		Spear:        "Spear",
		Spetum:       "Spetum",
		Spice:        "Spice",
		Statue:       "Statue",
		Stave:        "Stave",
		Stones:       "Stones",
		String:       "String",
		Sugar:        "Sugar",
		Sword:        "Sword",
		Tapestries:   "Tapestries",
		Tea:          "Tea",
		Tobacco:      "Tobacco",
		Trap:         "Trap",
		Trews:        "Trews",
		Trinket:      "Trinket",
		Trumpet:      "Trumpet",
		Urn:          "Urn",
		Wagons:       "Wagons",
		Wax:          "Wax",
	}
	// StringToEnum is a helper map for unmarshalling the enum
	StringToEnum = map[string]Item_e{
		"N/A":          None,
		"Adze":         Adze,
		"Arbalest":     Arbalest,
		"Arrows":       Arrows,
		"Axes":         Axes,
		"Backpack":     Backpack,
		"Ballistae":    Ballistae,
		"Bark":         Bark,
		"Barrel":       Barrel,
		"Bladder":      Bladder,
		"Blubber":      Blubber,
		"Boat":         Boat,
		"BoneArmour":   BoneArmour,
		"Bones":        Bones,
		"Bows":         Bows,
		"Bread":        Bread,
		"Breastplate":  Breastplate,
		"Candle":       Candle,
		"Canoes":       Canoes,
		"Carpets":      Carpets,
		"Catapult":     Catapult,
		"Cattle":       Cattle,
		"Cauldrons":    Cauldrons,
		"Chain":        Chain,
		"China":        China,
		"Clay":         Clay,
		"Cloth":        Cloth,
		"Clubs":        Clubs,
		"Coal":         Coal,
		"Coffee":       Coffee,
		"Coins":        Coins,
		"Cotton":       Cotton,
		"Cuirass":      Cuirass,
		"Cuirboilli":   Cuirboilli,
		"Diamond":      Diamond,
		"Diamonds":     Diamonds,
		"Drum":         Drum,
		"Elephant":     Elephant,
		"Falchion":     Falchion,
		"Fish":         Fish,
		"Flax":         Flax,
		"Flour":        Flour,
		"Flute":        Flute,
		"Fodder":       Fodder,
		"Frame":        Frame,
		"Frankincense": Frankincense,
		"Fur":          Fur,
		"Glasspipe":    Glasspipe,
		"Goats":        Goats,
		"Gold":         Gold,
		"Grain":        Grain,
		"Grape":        Grape,
		"Gut":          Gut,
		"HBow":         HBow,
		"Harp":         Harp,
		"Haube":        Haube,
		"Heaters":      Heaters,
		"Helm":         Helm,
		"Herbs":        Herbs,
		"Hive":         Hive,
		"Hoe":          Hoe,
		"Honey":        Honey,
		"Hood":         Hood,
		"Horn":         Horn,
		"Horses":       Horses,
		"Jade":         Jade,
		"Jerkin":       Jerkin,
		"Kayak":        Kayak,
		"Ladder":       Ladder,
		"Leather":      Leather,
		"Logs":         Logs,
		"Lute":         Lute,
		"Mace":         Mace,
		"Mattock":      Mattock,
		"Metal":        Metal,
		"MillStone":    MillStone,
		"Musk":         Musk,
		"Net":          Net,
		"Oar":          Oar,
		"Oil":          Oil,
		"Olives":       Olives,
		"Opium":        Opium,
		"Ores":         Ores,
		"Paddle":       Paddle,
		"Palanquin":    Palanquin,
		"Parchment":    Parchment,
		"Pavis":        Pavis,
		"Pearls":       Pearls,
		"Pellets":      Pellets,
		"People":       People,
		"Pewter":       Pewter,
		"Picks":        Picks,
		"Plows":        Plows,
		"Provisions":   Provisions,
		"Quarrel":      Quarrel,
		"Rake":         Rake,
		"Ram":          Ram,
		"Ramp":         Ramp,
		"Ring":         Ring,
		"Rope":         Rope,
		"Rug":          Rug,
		"Saddle":       Saddle,
		"Saddlebag":    Saddlebag,
		"Salt":         Salt,
		"Sand":         Sand,
		"Scale":        Scale,
		"Sculpture":    Sculpture,
		"Scutum":       Scutum,
		"Scythe":       Scythe,
		"Shackle":      Shackle,
		"Shaft":        Shaft,
		"Shield":       Shield,
		"Shovel":       Shovel,
		"Silk":         Silk,
		"Silver":       Silver,
		"Skin":         Skin,
		"Slaves":       Slaves,
		"Slings":       Slings,
		"Snare":        Snare,
		"Spear":        Spear,
		"Spetum":       Spetum,
		"Spice":        Spice,
		"Statue":       Statue,
		"Stave":        Stave,
		"Stones":       Stones,
		"String":       String,
		"Sugar":        Sugar,
		"Sword":        Sword,
		"Tapestries":   Tapestries,
		"Tea":          Tea,
		"Tobacco":      Tobacco,
		"Trap":         Trap,
		"Trews":        Trews,
		"Trinket":      Trinket,
		"Trumpet":      Trumpet,
		"Urn":          Urn,
		"Wagons":       Wagons,
		"Wax":          Wax,
	}

	LowerCaseName = map[string]Item_e{
		"adze":         Adze,
		"arbalest":     Arbalest,
		"arrows":       Arrows,
		"axes":         Axes,
		"backpack":     Backpack,
		"ballistae":    Ballistae,
		"bark":         Bark,
		"barrel":       Barrel,
		"bladder":      Bladder,
		"blubber":      Blubber,
		"boat":         Boat,
		"bonearmour":   BoneArmour,
		"bones":        Bones,
		"bows":         Bows,
		"bread":        Bread,
		"breastplate":  Breastplate,
		"candle":       Candle,
		"canoes":       Canoes,
		"carpets":      Carpets,
		"catapult":     Catapult,
		"cattle":       Cattle,
		"cauldrons":    Cauldrons,
		"chain":        Chain,
		"china":        China,
		"clay":         Clay,
		"cloth":        Cloth,
		"clubs":        Clubs,
		"coal":         Coal,
		"coffee":       Coffee,
		"coins":        Coins,
		"cotton":       Cotton,
		"cuirass":      Cuirass,
		"cuirboilli":   Cuirboilli,
		"diamond":      Diamond,
		"diamonds":     Diamonds,
		"drum":         Drum,
		"elephant":     Elephant,
		"falchion":     Falchion,
		"fish":         Fish,
		"flax":         Flax,
		"flour":        Flour,
		"flute":        Flute,
		"fodder":       Fodder,
		"frame":        Frame,
		"frankincense": Frankincense,
		"fur":          Fur,
		"glasspipe":    Glasspipe,
		"goats":        Goats,
		"gold":         Gold,
		"grain":        Grain,
		"grape":        Grape,
		"gut":          Gut,
		"hbow":         HBow,
		"harp":         Harp,
		"haube":        Haube,
		"heaters":      Heaters,
		"helm":         Helm,
		"herbs":        Herbs,
		"hive":         Hive,
		"hoe":          Hoe,
		"honey":        Honey,
		"hood":         Hood,
		"horn":         Horn,
		"horses":       Horses,
		"jade":         Jade,
		"jerkin":       Jerkin,
		"kayak":        Kayak,
		"ladder":       Ladder,
		"leather":      Leather,
		"logs":         Logs,
		"lute":         Lute,
		"mace":         Mace,
		"mattock":      Mattock,
		"metal":        Metal,
		"millstone":    MillStone,
		"musk":         Musk,
		"net":          Net,
		"oar":          Oar,
		"oil":          Oil,
		"olives":       Olives,
		"opium":        Opium,
		"ores":         Ores,
		"paddle":       Paddle,
		"palanquin":    Palanquin,
		"parchment":    Parchment,
		"pavis":        Pavis,
		"pearls":       Pearls,
		"pellets":      Pellets,
		"people":       People,
		"pewter":       Pewter,
		"picks":        Picks,
		"plows":        Plows,
		"provisions":   Provisions,
		"quarrel":      Quarrel,
		"rake":         Rake,
		"ram":          Ram,
		"ramp":         Ramp,
		"ring":         Ring,
		"rope":         Rope,
		"rug":          Rug,
		"saddle":       Saddle,
		"saddlebag":    Saddlebag,
		"salt":         Salt,
		"sand":         Sand,
		"scale":        Scale,
		"sculpture":    Sculpture,
		"scutum":       Scutum,
		"scythe":       Scythe,
		"shackle":      Shackle,
		"shaft":        Shaft,
		"shield":       Shield,
		"shovel":       Shovel,
		"silk":         Silk,
		"silver":       Silver,
		"skin":         Skin,
		"slaves":       Slaves,
		"slings":       Slings,
		"snare":        Snare,
		"spear":        Spear,
		"spetum":       Spetum,
		"spice":        Spice,
		"statue":       Statue,
		"stave":        Stave,
		"stones":       Stones,
		"string":       String,
		"sugar":        Sugar,
		"sword":        Sword,
		"tapestries":   Tapestries,
		"tea":          Tea,
		"tobacco":      Tobacco,
		"trap":         Trap,
		"trews":        Trews,
		"trinket":      Trinket,
		"trumpet":      Trumpet,
		"urn":          Urn,
		"wagons":       Wagons,
		"wax":          Wax,
	}
)
