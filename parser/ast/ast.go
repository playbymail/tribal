// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package ast

import (
	"github.com/playbymail/tribal/border"
	"github.com/playbymail/tribal/direction"
	"github.com/playbymail/tribal/passage"
	"github.com/playbymail/tribal/resource"
	"github.com/playbymail/tribal/terrain"
)

// ClanId_t defines a clan identifier;
// this is the unit that owns a turn report.
type ClanId_t string

// Unit_t defines a single section of a turn report
type Unit_t struct {
	Id          UnitId_t      `json:"id"`
	Name        UnitName_t    `json:"name,omitempty"`   // optional name field
	PreviousHex Coordinates_t `json:"previous_hex"`     // location of unit at the beginning of the turn
	CurrentHex  Coordinates_t `json:"current_hex"`      // location of unit at the end of the turn
	Turn        TurnId_t      `json:"turn,omitempty"`   // turn number
	Moves       *Moves_t      `json:"moves,omitempty"`  // moves made by the unit
	Status      *Status_t     `json:"status,omitempty"` // status of the unit
	Errors      []error       `json:"errors,omitempty"`
}

// UnitId_t defines a unit identifier in a turn report
type UnitId_t string

// UnitName_t defines a unit name in a turn report
type UnitName_t string

// TurnId_t is an alternative turn identifier from a turn report.
type TurnId_t int

// Turn_t defines the turn year and month from a turn report.
type Turn_t struct {
	Year  int
	Month int
}

// Moves_t defines a node containing a unit's movement and results in a turn report.
type Moves_t struct {
	Follows *Follows_t  `json:"follows,omitempty"`
	GoesTo  *GoesTo_t   `json:"goes_to,omitempty"`
	Marches []*March_t  `json:"marches,omitempty"`
	Patrols []*Patrol_t `json:"patrols,omitempty"`
	Errors  []error     `json:"errors,omitempty"`
}

// Follows_t defines the results for a follows line
type Follows_t struct {
	Id      UnitId_t      `json:"id"`
	Follows UnitId_t      `json:"follows"`
	From    Coordinates_t `json:"from,omitempty"`
	To      Coordinates_t `json:"to,omitempty"`
}

// GoesTo_t defines the results for a goes to line
type GoesTo_t struct {
	Id     UnitId_t      `json:"id"`
	From   Coordinates_t `json:"from,omitempty"`
	GoesTo Coordinates_t `json:"goes_to,omitempty"`
	To     Coordinates_t `json:"to,omitempty"`
}

// March_t defines the results of a single segment of a unit's land-based movement.
type March_t struct {
	Id        UnitId_t              `json:"id"`
	From      Coordinates_t         `json:"from,omitempty"`
	Direction direction.Direction_e `json:"direction"`
	To        Coordinates_t         `json:"to,omitempty"`
	Terrain   terrain.Terrain_e     `json:"terrain"`
	Neighbors []*Neighbor_t         `json:"neighbors,omitempty"`
	Borders   []*Border_t           `json:"borders,omitempty"`
	Passages  []*Passage_t          `json:"passages,omitempty"`
	HexName   *HexName_t            `json:"hex_name,omitempty"`
	Errors    *MarchErrors_t        `json:"errors,omitempty"`
}

// Patrol_t defines the results of a single segment of a scout's patrol.
type Patrol_t struct {
	Id         UnitId_t              `json:"id"`
	Patrol     int                   `json:"patrol"`
	From       Coordinates_t         `json:"from,omitempty"`
	Direction  direction.Direction_e `json:"direction"`
	To         Coordinates_t         `json:"to,omitempty"`
	Terrain    terrain.Terrain_e     `json:"terrain"`
	Neighbors  []*Neighbor_t         `json:"neighbors,omitempty"`
	Borders    []*Border_t           `json:"borders,omitempty"`
	Passages   []*Passage_t          `json:"passages,omitempty"`
	Resources  []resource.Resource_e `json:"resources,omitempty"`
	Encounters []UnitId_t            `json:"encounters,omitempty"`
	HexName    *HexName_t            `json:"hex_name,omitempty"`
	Errors     *PatrolErrors_t       `json:"errors,omitempty"`
}

// Neighbor_t defines the direction and type of neighboring tile in a turn report.
// Note that not all terrain types are observable as neighbors.
type Neighbor_t struct {
	Terrain   terrain.Terrain_e       `json:"terrain"`
	Direction []direction.Direction_e `json:"direction,omitempty"`
}

// Border_t defines the direction and type of neighboring tile in a turn report.
// Note that not all terrain types are observable as neighbors.
type Border_t struct {
	Border    border.Border_e         `json:"border"`
	Direction []direction.Direction_e `json:"direction,omitempty"`
}

// Passage_t defines the direction and type of passage in a turn report.
// Note that passages are not supposed to be one-way, but errors in the
// map database may result in one-way passages. Players are expected to
// inform the GM, who will update the map. That means we must not treat
// newly formed passages as errors.
type Passage_t struct {
	Passage   passage.Passage_e       `json:"passage"`
	Direction []direction.Direction_e `json:"direction,omitempty"`
}

// MarchErrors_t defines some common errors encountered while processing a segment in a turn report.
type MarchErrors_t struct {
	ExcessInput []string `json:"excess_input,omitempty"`
	Errors      []error  `json:"errors,omitempty"`
}

// PatrolErrors_t defines some common errors encountered while processing a segment in a turn report.
type PatrolErrors_t struct {
	ExcessInput []string `json:"excess_input,omitempty"`
	Errors      []error  `json:"errors,omitempty"`
}

// Status_t defines the status line of a unit in a turn report.
type Status_t struct {
	Unit   UnitId_t        `json:"unit,omitempty"`
	Tile   Tile_t          `json:"tile,omitempty"`
	Errors *StatusErrors_t `json:"errors,omitempty"`
}

type StatusErrors_t struct {
	ExcessInput []string `json:"excess_input,omitempty"`
	Errors      []error  `json:"errors,omitempty"`
}

type Tile_t struct {
	Coordinates Coordinates_t         `json:"coordinates"`
	Terrain     terrain.Terrain_e     `json:"terrain,omitempty"`
	HexName     *HexName_t            `json:"hex_name,omitempty"`
	Resources   []resource.Resource_e `json:"resources,omitempty"`
	Neighbors   []*Neighbor_t         `json:"neighbors,omitempty"`
	Borders     []*Border_t           `json:"borders,omitempty"`
	Passages    []*Passage_t          `json:"passages,omitempty"`
	Encounters  []UnitId_t            `json:"encounters,omitempty"`
}
