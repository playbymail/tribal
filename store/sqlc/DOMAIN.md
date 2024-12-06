# Data

NOTE:
Almost all tables will have Clan as part of the key. This allows us to store reports for multiple players while preventing them from seeing each other's data.

## Key Types

1. Turn - a turn is identified by a Year and a Month. It is always formatted as YYYY-MM.
2. Year - a year in the game is an integer in the range 899 through 1234.
3. Month - a month in the game is an integer in the range 1 through 12.
4. Clan - players are identified by Clan. It is an integer in the range 1 through 999. It is always formatted as a four-digit number.
5. Unit - a unit is an element that can move. It is a string consisting of the Clan followed by a two-character suffix.
6. Hex - a hex is a tile on the game map. Hexes on the map are flat-topped. They are identified by their row and column, which are both positive integers.
7. Tile - a hex on the generated map. Tiles have features, terrain, borders, and resources. They are identified on the map using grid, row, and column.
8. Direction - a code indicating the direction of movement. Examples are "N" (north), "SE" (southeast), and "S" (south).
9. Terrain - a code describing the type of terrain in a hex. The Terrain is used to generate the map. Examples are "PR" (prairie) or "SW" (swamp).
10. Border - a code describing the border separating two hexes. Examples are "R" (river) or "C" (canal).
11. Passage - a code describing a passage across or through the border separating two hexes. Examples are "SR" (stone road) or "MP" (mountain pass).
12. Resources - a code describing a resource in a hex. Examples are "Coal Mine" or "Iron Ore."
13. Settlement - a settlement is identified by its location and a name. Settlements are not permanent, so we must track the turns that it existed.
14. Transients - units are free to move between hexes. Transient records capture where they ended a turn.

## Effective Dating records
Most tables use two columns, EffDt and EndDt, to determine when the record is valid.
Both of these columns refer to turn dates, not calendar dates.

The EffDt column determines when a record becomes active.
The EndDt column determines when a record becomes inactive.
It is important to remember that a record is inactive starting on the EndDt.

Queries on these tables must use an "as of date" to select a record.

## Tile
A tile is a hex on the game board which is used when generating the map.
Tiles are uniquely identified by their grid, row, and column.
The engine that renders maps is responsible for converting the hex location to the tile location.

Tiles are assigned attributes such as terrain, borders, passages, and resources.
This data is added to the Tile as units explore the game board.

## Transients
Transients are identified by Turn, Location, and Unit.
Multiple units may occupy the same location during a turn.

## Settlements
Settlements are identified by a generated integer.
Settlements are built by players, maintained, and may be abandoned, so we must track their duration.
We do this with two columns, EffDt and EndDt.
EffDt is the first Turn that a Settlement was seen.
EffDt is initially set to the last turn (9999-12) and is only updated when the settlement is abandoned.
When that happens, EffDt is updated.
For mapping purposes, a Settlement is placed on the map starting on EffDt and removed starting on EndDt.
Settlement names are not required to be unique.

## Movement

Units may travel between hexes during the course of the game.
Move captures the starting location, the direction of travel, and all results of the move.
A Move is uniquely identified by the Unit, Turn, and a sequence number indicating the order of the move within the turn.

A Step captures the results of a Move.

1. Starting Location - the Hex the Unit was moving from.
2. Direction - the Direction the Unit wanted to move.
3. Succeeded - a boolean that is true only if the Move succeeded.
4. Reason For Failure - text describing why the Move failed. Optional and set only if the Move failed.
5. Ending Location - the Hex the Unit ended up in. If the Move failed, this will be the same value as the Starting Location.
6. Terrain - the Terrain of the Ending Location's hex.
7. Borders - each edge of the hex may have a border.
8. Passages - each edge of the hex may have a passage.
9. Settlement - a step may find a settlement.
10. Units - a step may find multiple units.

