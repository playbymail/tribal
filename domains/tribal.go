// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import (
	"fmt"
	"github.com/playbymail/tribal/direction"
	"github.com/playbymail/tribal/resource"
	"github.com/playbymail/tribal/terrain"
)

type Unit_t struct {
	Id          UnitId_t
	Name        UnitName_t    // optional name field
	PreviousHex Coordinates_t // location of unit at the beginning of the turn
	CurrentHex  Coordinates_t // location of unit at the end of the turn
	Turn        TurnId_t      // turn number
	Status      *Status_t     // status of the unit
	Error       error         // highest level error encountered while parsing the unit
}

type UnitId_t string

type UnitName_t string

type Coordinates_t struct {
	GridRow    int // 1-based, A ... Z -> 1 ... 26
	GridColumn int // 1-based, A ... Z -> 1 ... 26
	Column     int // 1-based, 1 ... 30
	Row        int // 1-based, 1 ... 21
}

func (c Coordinates_t) String() string {
	if c.GridRow == 0 && c.GridColumn == 0 && c.Column == 0 && c.Row == 0 {
		return "n/a"
	} else if c.GridRow == 0 && c.GridColumn == 0 {
		return fmt.Sprintf("## %02d%02d", c.Column, c.Row)
	}
	return fmt.Sprintf("%c%c %02d%02d", c.GridRow+'A'-1, c.GridColumn+'A'-1, c.Column, c.Row)
}

type Neighbor_t struct {
	Direction direction.Direction_e
	Terrain   terrain.Terrain_e
}

type Status_t struct {
	Unit           UnitId_t
	Terrain        terrain.Terrain_e
	SettlementName string
	Resources      resource.Resource_e
	Encounters     []UnitId_t
}

type Tile_t struct {
	Coordinates Coordinates_t
	Terrain     terrain.Terrain_e
	Neighbors   []Neighbor_t
	Resources   resource.Resource_e
	Units       []Unit_t
}

type Turn_t struct {
	Year  int
	Month int
}

type TurnId_t int
