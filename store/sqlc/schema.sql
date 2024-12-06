-- Define the tables for the schema

PRAGMA foreign_keys = OFF;
DROP TABLE IF EXISTS border_codes;
DROP TABLE IF EXISTS clans;
DROP TABLE IF EXISTS item_codes;
DROP TABLE IF EXISTS move_border_details;
DROP TABLE IF EXISTS move_passage_details;
DROP TABLE IF EXISTS move_resource_details;
DROP TABLE IF EXISTS move_settlement_details;
DROP TABLE IF EXISTS move_transient_details;
DROP TABLE IF EXISTS moves;
DROP TABLE IF EXISTS passage_codes;
DROP TABLE IF EXISTS report_files;
DROP TABLE IF EXISTS resource_codes;
DROP TABLE IF EXISTS terrain_codes;
DROP TABLE IF EXISTS tile_border_details;
DROP TABLE IF EXISTS tile_passage_details;
DROP TABLE IF EXISTS tile_resource_details;
DROP TABLE IF EXISTS tile_settlement_details;
DROP TABLE IF EXISTS tile_terrain_details;
DROP TABLE IF EXISTS tile_transient_details;
DROP TABLE IF EXISTS tiles;
DROP TABLE IF EXISTS turns;
DROP TABLE IF EXISTS units;
PRAGMA foreign_keys = ON;

-- --------------------------------------------------------------------------
-- Report Files
--
-- This table contains the hash and name of reports files that we've loaded.
-- This is used to prevent duplicate reports from being loaded.
-- The database is single-user, so we don't need to worry about multiple
-- players uploading the same file.
--
-- The hash is the SHA-1 hash of the report.
--
-- We don't care about the name of the file, so there are no constraints
-- on it. We store it so that players can see what they've loaded based
-- on the file name on their computer.
CREATE TABLE report_files
(
    id INTEGER NOT NULL PRIMARY KEY,
    hash TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (CAST(strftime('%s', 'now') AS INTEGER))
);

-- --------------------------------------------------------------------------
-- Turns
--
-- The application assumes that we have a start turn of 899-12 and
-- end turn of 9999-12 pre-populated.
--
-- Note: the formula "((year - 899) * 12) + month - 12" is used to
-- convert the year and month into the TribeNet turn number, which
-- starts at 0 for turn 899-12.
CREATE TABLE turns
(
    id    INTEGER NOT NULL PRIMARY KEY, -- calculated as (year-899) * 12 + month - 12
    year  INTEGER NOT NULL CHECK (year BETWEEN 899 AND 9999),
    month INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12),
    UNIQUE (year, month)
);

INSERT INTO turns (id, year, month)
VALUES (0, 899, 12);
INSERT INTO turns (id, year, month)
VALUES (1, 900, 1);
INSERT INTO turns (id, year, month)
VALUES (2, 900, 2);
INSERT INTO turns (id, year, month)
VALUES (3, 900, 3);
INSERT INTO turns (id, year, month)
VALUES (4, 900, 4);
INSERT INTO turns (id, year, month)
VALUES (5, 900, 5);
INSERT INTO turns (id, year, month)
VALUES (6, 900, 6);
INSERT INTO turns (id, year, month)
VALUES (7, 900, 7);
INSERT INTO turns (id, year, month)
VALUES (8, 900, 8);
INSERT INTO turns (id, year, month)
VALUES (9, 900, 9);
INSERT INTO turns (id, year, month)
VALUES (10, 900, 10);
INSERT INTO turns (id, year, month)
VALUES (11, 900, 11);
INSERT INTO turns (id, year, month)
VALUES (12, 900, 12);
INSERT INTO turns (id, year, month)
VALUES (13, 901, 1);
INSERT INTO turns (id, year, month)
VALUES (14, 901, 2);
INSERT INTO turns (id, year, month)
VALUES (15, 901, 3);
INSERT INTO turns (id, year, month)
VALUES (16, 901, 4);
INSERT INTO turns (id, year, month)
VALUES (17, 901, 5);
INSERT INTO turns (id, year, month)
VALUES (18, 901, 6);
INSERT INTO turns (id, year, month)
VALUES (19, 901, 7);
INSERT INTO turns (id, year, month)
VALUES (20, 901, 8);
INSERT INTO turns (id, year, month)
VALUES (21, 901, 9);
INSERT INTO turns (id, year, month)
VALUES (22, 901, 10);
INSERT INTO turns (id, year, month)
VALUES (23, 901, 11);
INSERT INTO turns (id, year, month)
VALUES (24, 901, 12);
INSERT INTO turns (id, year, month)
VALUES (25, 902, 1);
INSERT INTO turns (id, year, month)
VALUES (26, 902, 2);
INSERT INTO turns (id, year, month)
VALUES (27, 902, 3);
INSERT INTO turns (id, year, month)
VALUES (28, 902, 4);
INSERT INTO turns (id, year, month)
VALUES (29, 902, 5);
INSERT INTO turns (id, year, month)
VALUES (30, 902, 6);
INSERT INTO turns (id, year, month)
VALUES (31, 902, 7);
INSERT INTO turns (id, year, month)
VALUES (32, 902, 8);
INSERT INTO turns (id, year, month)
VALUES (33, 902, 9);
INSERT INTO turns (id, year, month)
VALUES (34, 902, 10);
INSERT INTO turns (id, year, month)
VALUES (35, 902, 11);
INSERT INTO turns (id, year, month)
VALUES (36, 902, 12);
INSERT INTO turns (id, year, month)
VALUES ((9999 - 899) * 12 + 12 - 12, 9999, 12);

-- --------------------------------------------------------------------------
-- Clans
--
-- We currently store only one clan per database for security,
-- so this field is for reference only until that changes
-- (which will be never since it requires buy-in from all players).
CREATE TABLE clans
(
    id   INTEGER NOT NULL PRIMARY KEY CHECK (id BETWEEN 1 AND 999),
    name TEXT    NOT NULL, -- the formatted name of the clan, e.g. 0987
    UNIQUE (name)
);

-- --------------------------------------------------------------------------
-- Units
--
-- Note that we never want to add scout units to the Transients table.
-- The format of the id is xxxx for the clan and tribes,
-- xxxx([cefg][1-9]) for couriers, elements, fleets, and garrisons,
-- and xxxx([cefg][1-9])?(s[1-8]) for scouts.
CREATE TABLE units
(
    id       TEXT PRIMARY KEY,
    clan_no  INTEGER NOT NULL REFERENCES clans (id),               -- not really needed for single-clan database
    is_scout INTEGER NOT NULL DEFAULT 0 CHECK (is_scout in (0, 1)) -- true only if unit is a scout
);

-- --------------------------------------------------------------------------
-- Border Codes
--
-- This table stores the codes that describe a tile border.
CREATE TABLE border_codes
(
    code        TEXT NOT NULL PRIMARY KEY, -- R, CANAL, etc.
    descr       TEXT NOT NULL,             -- River, Canal, etc.
    wxx_feature TEXT NOT NULL,             -- how to draw the border in worldographer
    UNIQUE (descr)
);

INSERT INTO border_codes
VALUES ('CANAL', 'Canal', '*');
INSERT INTO border_codes
VALUES ('RIVER', 'River', '*');

-- --------------------------------------------------------------------------
-- Item Codes
--
-- This table stores the codes that describe an item that can be found in a tile.
--
-- I don't think we need to track these; the map render is not even aware of them.
CREATE TABLE item_codes
(
    code  TEXT NOT NULL PRIMARY KEY, -- JEWELS, PONIES, RICH PERSON, etc.
    descr TEXT NOT NULL,             -- Jewels, Ponies, Rich Person, etc.
    UNIQUE (descr)
);

INSERT INTO item_codes
VALUES ('ADZE', 'Adze');
INSERT INTO item_codes
VALUES ('ARBALEST', 'Arbalest');
INSERT INTO item_codes
VALUES ('ARROWS', 'Arrows');
INSERT INTO item_codes
VALUES ('AXES', 'Axes');
INSERT INTO item_codes
VALUES ('BACKPACK', 'Backpack');
INSERT INTO item_codes
VALUES ('BALLISTAE', 'Ballistae');
INSERT INTO item_codes
VALUES ('BARK', 'Bark');
INSERT INTO item_codes
VALUES ('BARREL', 'Barrel');
INSERT INTO item_codes
VALUES ('BLADDER', 'Bladder');
INSERT INTO item_codes
VALUES ('BLUBBER', 'Blubber');
INSERT INTO item_codes
VALUES ('BOAT', 'Boat');
INSERT INTO item_codes
VALUES ('BONEARMOUR', 'BoneArmour');
INSERT INTO item_codes
VALUES ('BONES', 'Bones');
INSERT INTO item_codes
VALUES ('BOWS', 'Bows');
INSERT INTO item_codes
VALUES ('BREAD', 'Bread');
INSERT INTO item_codes
VALUES ('BREASTPLATE', 'Breastplate');
INSERT INTO item_codes
VALUES ('CANDLE', 'Candle');
INSERT INTO item_codes
VALUES ('CANOES', 'Canoes');
INSERT INTO item_codes
VALUES ('CARPETS', 'Carpets');
INSERT INTO item_codes
VALUES ('CATAPULT', 'Catapult');
INSERT INTO item_codes
VALUES ('CATTLE', 'Cattle');
INSERT INTO item_codes
VALUES ('CAULDRONS', 'Cauldrons');
INSERT INTO item_codes
VALUES ('CHAIN', 'Chain');
INSERT INTO item_codes
VALUES ('CHINA', 'China');
INSERT INTO item_codes
VALUES ('CLAY', 'Clay');
INSERT INTO item_codes
VALUES ('CLOTH', 'Cloth');
INSERT INTO item_codes
VALUES ('CLUBS', 'Clubs');
INSERT INTO item_codes
VALUES ('COAL', 'Coal');
INSERT INTO item_codes
VALUES ('COFFEE', 'Coffee');
INSERT INTO item_codes
VALUES ('COINS', 'Coins');
INSERT INTO item_codes
VALUES ('COTTON', 'Cotton');
INSERT INTO item_codes
VALUES ('CUIRASS', 'Cuirass');
INSERT INTO item_codes
VALUES ('CUIRBOILLI', 'Cuirboilli');
INSERT INTO item_codes
VALUES ('DIAMOND', 'Diamond');
INSERT INTO item_codes
VALUES ('DIAMONDS', 'Diamonds');
INSERT INTO item_codes
VALUES ('DRUM', 'Drum');
INSERT INTO item_codes
VALUES ('ELEPHANT', 'Elephant');
INSERT INTO item_codes
VALUES ('FALCHION', 'Falchion');
INSERT INTO item_codes
VALUES ('FISH', 'Fish');
INSERT INTO item_codes
VALUES ('FLAX', 'Flax');
INSERT INTO item_codes
VALUES ('FLOUR', 'Flour');
INSERT INTO item_codes
VALUES ('FLUTE', 'Flute');
INSERT INTO item_codes
VALUES ('FODDER', 'Fodder');
INSERT INTO item_codes
VALUES ('FRAME', 'Frame');
INSERT INTO item_codes
VALUES ('FRANKINCENSE', 'Frankincense');
INSERT INTO item_codes
VALUES ('FUR', 'Fur');
INSERT INTO item_codes
VALUES ('GLASSPIPE', 'Glasspipe');
INSERT INTO item_codes
VALUES ('GOATS', 'Goats');
INSERT INTO item_codes
VALUES ('GOLD', 'Gold');
INSERT INTO item_codes
VALUES ('GRAIN', 'Grain');
INSERT INTO item_codes
VALUES ('GRAPE', 'Grape');
INSERT INTO item_codes
VALUES ('GUT', 'Gut');
INSERT INTO item_codes
VALUES ('HBOW', 'HBow');
INSERT INTO item_codes
VALUES ('HARP', 'Harp');
INSERT INTO item_codes
VALUES ('HAUBE', 'Haube');
INSERT INTO item_codes
VALUES ('HEATERS', 'Heaters');
INSERT INTO item_codes
VALUES ('HELM', 'Helm');
INSERT INTO item_codes
VALUES ('HERBS', 'Herbs');
INSERT INTO item_codes
VALUES ('HIVE', 'Hive');
INSERT INTO item_codes
VALUES ('HOE', 'Hoe');
INSERT INTO item_codes
VALUES ('HONEY', 'Honey');
INSERT INTO item_codes
VALUES ('HOOD', 'Hood');
INSERT INTO item_codes
VALUES ('HORN', 'Horn');
INSERT INTO item_codes
VALUES ('HORSES', 'Horses');
INSERT INTO item_codes
VALUES ('JADE', 'Jade');
INSERT INTO item_codes
VALUES ('JERKIN', 'Jerkin');
INSERT INTO item_codes
VALUES ('KAYAK', 'Kayak');
INSERT INTO item_codes
VALUES ('LADDER', 'Ladder');
INSERT INTO item_codes
VALUES ('LEATHER', 'Leather');
INSERT INTO item_codes
VALUES ('LOGS', 'Logs');
INSERT INTO item_codes
VALUES ('LUTE', 'Lute');
INSERT INTO item_codes
VALUES ('MACE', 'Mace');
INSERT INTO item_codes
VALUES ('MATTOCK', 'Mattock');
INSERT INTO item_codes
VALUES ('METAL', 'Metal');
INSERT INTO item_codes
VALUES ('MILLSTONE', 'MillStone');
INSERT INTO item_codes
VALUES ('MUSK', 'Musk');
INSERT INTO item_codes
VALUES ('NET', 'Net');
INSERT INTO item_codes
VALUES ('OAR', 'Oar');
INSERT INTO item_codes
VALUES ('OIL', 'Oil');
INSERT INTO item_codes
VALUES ('OLIVES', 'Olives');
INSERT INTO item_codes
VALUES ('OPIUM', 'Opium');
INSERT INTO item_codes
VALUES ('ORES', 'Ores');
INSERT INTO item_codes
VALUES ('PADDLE', 'Paddle');
INSERT INTO item_codes
VALUES ('PALANQUIN', 'Palanquin');
INSERT INTO item_codes
VALUES ('PARCHMENT', 'Parchment');
INSERT INTO item_codes
VALUES ('PAVIS', 'Pavis');
INSERT INTO item_codes
VALUES ('PEARLS', 'Pearls');
INSERT INTO item_codes
VALUES ('PELLETS', 'Pellets');
INSERT INTO item_codes
VALUES ('PEOPLE', 'People');
INSERT INTO item_codes
VALUES ('PEWTER', 'Pewter');
INSERT INTO item_codes
VALUES ('PICKS', 'Picks');
INSERT INTO item_codes
VALUES ('PLOWS', 'Plows');
INSERT INTO item_codes
VALUES ('PROVISIONS', 'Provisions');
INSERT INTO item_codes
VALUES ('QUARREL', 'Quarrel');
INSERT INTO item_codes
VALUES ('RAKE', 'Rake');
INSERT INTO item_codes
VALUES ('RAM', 'Ram');
INSERT INTO item_codes
VALUES ('RAMP', 'Ramp');
INSERT INTO item_codes
VALUES ('RING', 'Ring');
INSERT INTO item_codes
VALUES ('ROPE', 'Rope');
INSERT INTO item_codes
VALUES ('RUG', 'Rug');
INSERT INTO item_codes
VALUES ('SADDLE', 'Saddle');
INSERT INTO item_codes
VALUES ('SADDLEBAG', 'Saddlebag');
INSERT INTO item_codes
VALUES ('SALT', 'Salt');
INSERT INTO item_codes
VALUES ('SAND', 'Sand');
INSERT INTO item_codes
VALUES ('SCALE', 'Scale');
INSERT INTO item_codes
VALUES ('SCULPTURE', 'Sculpture');
INSERT INTO item_codes
VALUES ('SCUTUM', 'Scutum');
INSERT INTO item_codes
VALUES ('SCYTHE', 'Scythe');
INSERT INTO item_codes
VALUES ('SHACKLE', 'Shackle');
INSERT INTO item_codes
VALUES ('SHAFT', 'Shaft');
INSERT INTO item_codes
VALUES ('SHIELD', 'Shield');
INSERT INTO item_codes
VALUES ('SHOVEL', 'Shovel');
INSERT INTO item_codes
VALUES ('SILK', 'Silk');
INSERT INTO item_codes
VALUES ('SILVER', 'Silver');
INSERT INTO item_codes
VALUES ('SKIN', 'Skin');
INSERT INTO item_codes
VALUES ('SLAVES', 'Slaves');
INSERT INTO item_codes
VALUES ('SLINGS', 'Slings');
INSERT INTO item_codes
VALUES ('SNARE', 'Snare');
INSERT INTO item_codes
VALUES ('SPEAR', 'Spear');
INSERT INTO item_codes
VALUES ('SPETUM', 'Spetum');
INSERT INTO item_codes
VALUES ('SPICE', 'Spice');
INSERT INTO item_codes
VALUES ('STATUE', 'Statue');
INSERT INTO item_codes
VALUES ('STAVE', 'Stave');
INSERT INTO item_codes
VALUES ('STONES', 'Stones');
INSERT INTO item_codes
VALUES ('STRING', 'String');
INSERT INTO item_codes
VALUES ('SUGAR', 'Sugar');
INSERT INTO item_codes
VALUES ('SWORD', 'Sword');
INSERT INTO item_codes
VALUES ('TAPESTRIES', 'Tapestries');
INSERT INTO item_codes
VALUES ('TEA', 'Tea');
INSERT INTO item_codes
VALUES ('TOBACCO', 'Tobacco');
INSERT INTO item_codes
VALUES ('TRAP', 'Trap');
INSERT INTO item_codes
VALUES ('TREWS', 'Trews');
INSERT INTO item_codes
VALUES ('TRINKET', 'Trinket');
INSERT INTO item_codes
VALUES ('TRUMPET', 'Trumpet');
INSERT INTO item_codes
VALUES ('URN', 'Urn');
INSERT INTO item_codes
VALUES ('WAGONS', 'Wagons');
INSERT INTO item_codes
VALUES ('WAX', 'Wax');

-- --------------------------------------------------------------------------
-- Passage Codes
--
-- This table stores the codes that describe a tile passage.
CREATE TABLE passage_codes
(
    code        TEXT NOT NULL PRIMARY KEY, -- FORD, PASS, STONY ROAD, etc.
    descr       TEXT NOT NULL,             -- Ford, Mountain Pass, Stony Road, etc.
    wxx_feature TEXT NOT NULL,             -- how to draw the passage in worldographer
    UNIQUE (descr)
);

INSERT INTO passage_codes
VALUES ('FORD', 'Ford', '*');
INSERT INTO passage_codes
VALUES ('PASS', 'Pass', '*');
INSERT INTO passage_codes
VALUES ('STONEROAD', 'Stone Road', '*');

-- --------------------------------------------------------------------------
-- Resource Codes
--
-- This table stores the codes that describe a tile resource.
CREATE TABLE resource_codes
(
    code        TEXT NOT NULL PRIMARY KEY, -- COAL, IRON ORE, etc.
    descr       TEXT NOT NULL,             -- Coal, Iron Orer, etc.
    wxx_feature TEXT NOT NULL,             -- how to draw the resource in worldographer
    UNIQUE (descr)
);

INSERT INTO resource_codes
VALUES ('COAL', 'Coal', '*');
INSERT INTO resource_codes
VALUES ('COPPERORE', 'Copper Ore', '*');
INSERT INTO resource_codes
VALUES ('DIAMOND', 'Diamond', '*');
INSERT INTO resource_codes
VALUES ('FRANKINCENSE', 'Frankincense', '*');
INSERT INTO resource_codes
VALUES ('GOLD', 'Gold', '*');
INSERT INTO resource_codes
VALUES ('IRONORE', 'Iron Ore', '*');
INSERT INTO resource_codes
VALUES ('JADE', 'Jade', '*');
INSERT INTO resource_codes
VALUES ('KAOLIN', 'Kaolin', '*');
INSERT INTO resource_codes
VALUES ('LEADORE', 'Lead Ore', '*');
INSERT INTO resource_codes
VALUES ('LIMESTONE', 'Limestone', '*');
INSERT INTO resource_codes
VALUES ('NICKELORE', 'Nickel Ore', '*');
INSERT INTO resource_codes
VALUES ('PEARLS', 'Pearls', '*');
INSERT INTO resource_codes
VALUES ('PYRITE', 'Pyrite', '*');
INSERT INTO resource_codes
VALUES ('RUBIES', 'Rubies', '*');
INSERT INTO resource_codes
VALUES ('SALT', 'Salt', '*');
INSERT INTO resource_codes
VALUES ('SILVER', 'Silver', '*');
INSERT INTO resource_codes
VALUES ('SULPHUR', 'Sulphur', '*');
INSERT INTO resource_codes
VALUES ('TINORE', 'Tin Ore', '*');
INSERT INTO resource_codes
VALUES ('VANADIUMORE', 'Vanadium Ore', '*');
INSERT INTO resource_codes
VALUES ('ZINCORE', 'Zinc Ore', '*');

-- --------------------------------------------------------------------------
-- Terrain Codes
--
-- This table stores the codes that describe a tile terrain.
CREATE TABLE terrain_codes
(
    code        TEXT    NOT NULL PRIMARY KEY, -- PR, LJM, etc
    is_hills    INTEGER NOT NULL DEFAULT 0 CHECK (is_hills IN (0, 1)),
    is_jungle   INTEGER NOT NULL DEFAULT 0 CHECK (is_jungle IN (0, 1)),
    is_mountain INTEGER NOT NULL DEFAULT 0 CHECK (is_mountain IN (0, 1)),
    is_swamp    INTEGER NOT NULL DEFAULT 0 CHECK (is_swamp IN (0, 1)),
    is_water    INTEGER NOT NULL DEFAULT 0 CHECK (is_water IN (0, 1)),
    long_code   TEXT    NOT NULL,             -- PRAIRIE, LOW JUNGLE MOUNTAINS, etc
    descr       TEXT    NOT NULL,
    wxx_terrain TEXT    NOT NULL,
    UNIQUE (long_code),
    UNIQUE (descr)
);

INSERT INTO terrain_codes
VALUES ('*', 0, 0, 0, 0, 0, 'BLANK', 'Blank', '*');
INSERT INTO terrain_codes
VALUES ('ALPS', 0, 0, 1, 0, 0, 'ALPS', 'Alps', '*');
INSERT INTO terrain_codes
VALUES ('AH', 1, 0, 0, 0, 0, 'ARID_HILLS', 'Arid Hills', '*');
INSERT INTO terrain_codes
VALUES ('AR', 0, 0, 0, 0, 0, 'ARID_TUNDRA', 'Arid Tundra', '*');
INSERT INTO terrain_codes
VALUES ('BF', 0, 0, 0, 0, 0, 'BRUSH_FLAT', 'Brush Flat', '*');
INSERT INTO terrain_codes
VALUES ('BH', 1, 0, 0, 0, 0, 'BRUSH_HILLS', 'Brush Hills', '*');
INSERT INTO terrain_codes
VALUES ('CH', 1, 0, 0, 0, 0, 'CONIFER_HILLS', 'Conifer Hills', '*');
INSERT INTO terrain_codes
VALUES ('D', 0, 0, 0, 0, 0, 'DECIDUOUS', 'Deciduous', '*');
INSERT INTO terrain_codes
VALUES ('DE', 1, 0, 0, 0, 0, 'DECIDUOUS_HILLS', 'Deciduous Hills', '*');
INSERT INTO terrain_codes
VALUES ('DH', 0, 0, 0, 0, 0, 'DESERT', 'Desert', '*');
INSERT INTO terrain_codes
VALUES ('GH', 1, 0, 0, 0, 0, 'GRASSY_HILLS', 'Grassy Hills', '*');
INSERT INTO terrain_codes
VALUES ('GHP', 1, 0, 0, 0, 0, 'GRASSY_HILLS_PLATEAU', 'Grassy Hills Plateau', '*');
INSERT INTO terrain_codes
VALUES ('HSM', 0, 0, 1, 0, 0, 'HIGH_SNOWY_MOUNTAINS', 'High Snowy Mountains', '*');
INSERT INTO terrain_codes
VALUES ('JG', 0, 1, 0, 0, 0, 'JUNGLE', 'Jungle', '*');
INSERT INTO terrain_codes
VALUES ('JH', 1, 1, 0, 0, 0, 'JUNGLE_HILLS', 'Jungle Hills', '*');
INSERT INTO terrain_codes
VALUES ('L', 0, 0, 0, 0, 1, 'LAKE', 'Lake', '*');
INSERT INTO terrain_codes
VALUES ('LAM', 0, 0, 1, 0, 0, 'LOW_ARID_MOUNTAINS', 'Low Arid Mountains', '*');
INSERT INTO terrain_codes
VALUES ('LCM', 0, 0, 1, 0, 0, 'LOW_CONIFER_MOUNTAINS', 'Low Conifer Mountains', '*');
INSERT INTO terrain_codes
VALUES ('LJM', 0, 0, 1, 0, 0, 'LOW_JUNGLE_MOUNTAINS', 'Low Jungle Mountains', '*');
INSERT INTO terrain_codes
VALUES ('LSM', 0, 0, 1, 0, 0, 'LOW_SNOWY_MOUNTAINS', 'Low Snowy Mountains', '*');
INSERT INTO terrain_codes
VALUES ('LVM', 0, 0, 1, 0, 0, 'LOW_VOLCANIC_MOUNTAINS', 'Low Volcanic Mountains', '*');
INSERT INTO terrain_codes
VALUES ('O', 0, 0, 0, 0, 1, 'OCEAN', 'Ocean', '*');
INSERT INTO terrain_codes
VALUES ('PI', 0, 0, 0, 0, 0, 'POLAR_ICE', 'Polar Ice', '*');
INSERT INTO terrain_codes
VALUES ('PR', 0, 0, 0, 0, 0, 'PRAIRIE', 'Prairie', '*');
INSERT INTO terrain_codes
VALUES ('PPR', 0, 0, 0, 0, 0, 'PRAIRIE_PLATEAU', 'Prairie Plateau', '*');
INSERT INTO terrain_codes
VALUES ('RH', 1, 0, 0, 0, 0, 'ROCKY_HILLS', 'Rocky Hills', '*');
INSERT INTO terrain_codes
VALUES ('SH', 1, 0, 0, 0, 0, 'SNOWY_HILLS', 'Snowy Hills', '*');
INSERT INTO terrain_codes
VALUES ('SW', 0, 0, 0, 1, 0, 'SWAMP', 'Swamp', '*');
INSERT INTO terrain_codes
VALUES ('TU', 0, 0, 0, 0, 0, 'TUNDRA', 'Tundra', '*');
INSERT INTO terrain_codes
VALUES ('UJS', 0, 0, 0, 0, 0, 'UNKNOWN_JUNGLE_SWAMP', 'Unknown Jungle Swamp', '*');
INSERT INTO terrain_codes
VALUES ('UL', 0, 0, 0, 0, 0, 'UNKNOWN_LAND', 'Unknown Land', '*');
INSERT INTO terrain_codes
VALUES ('UM', 0, 0, 0, 0, 0, 'UNKNOWN_MOUNTAIN', 'Unknown Mountain', '*');
INSERT INTO terrain_codes
VALUES ('UW', 0, 0, 0, 0, 0, 'UNKNOWN_WATER', 'Unknown Water', '*');

-- --------------------------------------------------------------------------
-- the tile tables are used to render the map. the map generator understands
-- the effective date logic on the tables and uses it to create maps that
-- show the results "as of" a particular turn. future generators might even
-- use that information to trace movement paths for units.
-- --------------------------------------------------------------------------

-- --------------------------------------------------------------------------
-- Tiles
--
-- The direction columns (north, south, etc.) link tiles to their neighbors.
-- I am not sure that they are needed, but they make navigation queries simpler.
--
-- The last visited/last scouted values can be derived from either the Moves or
-- Transients tables. They may be removed if they make updates too expensive.
--
-- Tile attributes are stored in child tables because the values can change from
-- turn to turn or even move to move. For example, Fleet Movement could report a
-- tile as Unknown Water in one move, and then as Ocean in another.
--
-- It would be great if we could build a unique key on grid, row, and col, but
-- we can't since the early turn report obscured the grid. This setup, though,
-- allows us to easily update the grid, row, and col when we are able to compute
-- their values.
--
-- Anyway, we have to treat them as mutable since players are required to provide
-- missing values for early turn reports. This values will likely be updated once
-- the player gets reports that have the actual grid values.
CREATE TABLE tiles
(
    id              INTEGER PRIMARY KEY,
    grid            TEXT    NOT NULL,              -- usually ## or AA through ZZ, sometimes N/A
    row             INTEGER NOT NULL,              -- 0 only when grid is N/A
    col             INTEGER NOT NULL,              -- 0 only when grid is N/A
    north           INTEGER REFERENCES tiles (id),
    north_east      INTEGER REFERENCES tiles (id),
    north_west      INTEGER REFERENCES tiles (id),
    south           INTEGER REFERENCES tiles (id),
    south_east      INTEGER REFERENCES tiles (id),
    south_west      INTEGER REFERENCES tiles (id),
    last_visited_on INTEGER REFERENCES turns (id), -- last turn the tile was visited by a unit
    last_scouted_on INTEGER REFERENCES turns (id)  -- last turn the tile was scouted by a unit
);

-- --------------------------------------------------------------------------
-- Tile Border Details
--
-- These are derived after parsing all the movement results for a turn.
-- In other words, these are the details for the tile at the end of the turn.
--
-- We have to treat tile borders as mutable data because there are bugs
-- in the report generation process.
--
-- Assumption: each border of a tile can contain only one border feature.
-- This is likely invalid because of bugs.
--
-- The application is responsible for ensuring that the effective dated logic remains
-- consistent for all rows.
CREATE TABLE tile_border_details
(
    tile_id   INTEGER NOT NULL REFERENCES tiles (id),
    effdt     INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes active
    enddt     INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes inactive
    border_cd TEXT    NOT NULL REFERENCES border_codes (code),
    direction TEXT    NOT NULL CHECK (direction in ('N', 'NE', 'SE', 'S', 'SW', 'NW')),
    PRIMARY KEY (tile_id, border_cd, direction, effdt)
);

-- --------------------------------------------------------------------------
-- Tile Passage Details
--
-- These are derived after parsing all the movement results for a turn.
-- In other words, these are the details for the tile at the end of the turn.
--
-- We have to treat tile passages as mutable data because there are bugs
-- in the report generator and parser.
--
-- Assumption: each border of a tile can contain only one border feature.
-- This is likely invalid because of bugs.
--
-- The application is responsible for ensuring that the effective dated logic remains
-- consistent for all rows.
CREATE TABLE tile_passage_details
(
    tile_id    INTEGER NOT NULL REFERENCES tiles (id),
    effdt      INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes active
    enddt      INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes inactive
    passage_cd TEXT    NOT NULL REFERENCES passage_codes (code),
    direction  TEXT    NOT NULL CHECK (direction in ('N', 'NE', 'SE', 'S', 'SW', 'NW')),
    PRIMARY KEY (tile_id, effdt, passage_cd, direction)
);

-- --------------------------------------------------------------------------
-- Tile Resource Details
--
-- These are derived after parsing all the movement results for a turn.
-- In other words, these are the details for the tile at the end of the turn.
--
-- We have to treat tile resources as mutable data because there are bugs
-- in the report generator and parser.
--
-- Assumption: each tile can contain only one resource.
-- This is something that should be verified (but might be invalid because
-- of bugs, anyway).
--
-- The application is responsible for ensuring that the effective dated logic remains
-- consistent for all rows.
CREATE TABLE tile_resource_details
(
    tile_id     INTEGER NOT NULL REFERENCES tiles (id),
    effdt       INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes active
    enddt       INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes inactive,
    resource_cd TEXT    NOT NULL REFERENCES resource_codes (code),
    PRIMARY KEY (tile_id, effdt, resource_cd)
);

-- --------------------------------------------------------------------------
-- Tile Settlement Details
--
-- These are derived after parsing all the movement results for a turn.
-- In other words, these are the details for the tile at the end of the turn.
--
-- We have to treat settlements as mutable data because they can be destroyed or
-- abandoned. Also, there are bugs in the report generator and parser.
--
-- Assumption: tiles shouldn't have multiple settlements but there
-- are bugs in the report generation process and the parser, so we
-- have to allow them. We will silently merge duplicate names into
-- a single row, though.
--
-- Known issue: players won't know that a settlement has been
-- abandoned or destroyed until they send a unit to its location.
--
-- The application is responsible for ensuring that the effective dated logic remains
-- consistent for all rows.
CREATE TABLE tile_settlement_details
(
    tile_id INTEGER NOT NULL REFERENCES tiles (id),
    effdt   INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes active
    enddt   INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes inactive,
    name    TEXT    NOT NULL,
    PRIMARY KEY (tile_id, effdt, name)
);

-- --------------------------------------------------------------------------
-- Tile Terrain Details
--
-- These are derived after parsing all the movement results for a turn.
-- In other words, these are the details for the tile at the end of the turn.
--
-- We have to treat tile terrain as mutable data because of Fleet Movement reports.
-- Also, there are bugs in the report generator and parser.
--
-- Assumption: each tile can contain multiple terrain codes because of Fleet Movement
-- reports and bugs.
--
-- The application is responsible for ensuring that the effective dated logic remains
-- consistent for all rows.
CREATE TABLE tile_terrain_details
(
    tile_id    INTEGER NOT NULL REFERENCES tiles (id),
    effdt      INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes active
    enddt      INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes inactive,
    terrain_cd TEXT    NOT NULL REFERENCES units (id),
    PRIMARY KEY (tile_id, effdt, terrain_cd)
);

-- --------------------------------------------------------------------------
-- Tile Transient Details
--
-- These are derived after parsing all the movement results for a turn.
-- In other words, these are the details for the tile at the end of the turn.
--
-- We have to treat tile transients as mutable data because units are mobile and
-- there are bugs in the report generator and parser.
--
-- Unintended benefit of this table is it tracks where every unit ends the turn
-- as well as the last known location for any unit. It might be useful to add an
-- attribute to track the turn the unit was last seen.
--
-- Note: we must not add scout units to this table. If I knew how to enforce that
-- with a check constraint, I would.
--
-- The application is responsible for ensuring that the effective dated logic remains
-- consistent for all rows.
CREATE TABLE tile_transient_details
(
    tile_id INTEGER NOT NULL REFERENCES tiles (id),
    effdt   INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes active
    enddt   INTEGER NOT NULL REFERENCES turns (id), -- turn the entry becomes inactive,
    unit_id TEXT    NOT NULL REFERENCES units (id),
    PRIMARY KEY (tile_id, effdt, unit_id)
);

-- --------------------------------------------------------------------------
-- the moves and move detail tables capture the results of each move.
-- they are not needed after the tile detail tables are updated.
-- it may be cheaper just to keep them in memory and not load them at all.
-- --------------------------------------------------------------------------

-- --------------------------------------------------------------------------
-- Moves
--
-- This table stores information on all of the moves paresed from the turn reports.
-- We're assuming that there's no need to track the entry back to the source.
--
-- If a move fails, starting_tile and ending_tile must be set to the same value.
--
-- Warning: The Follow and Goes To moves don't have directions.
--
-- We could use a synthetic key (turn + unit + step) but that would make querying
-- the child tables irksome.
--
-- TODO: Fleet Moves have to be integrated into this somehow.
CREATE TABLE moves
(
    id             INTEGER PRIMARY KEY, -- unique identifier for the movement
    clan_no        INTEGER NOT NULL REFERENCES clans (id),
    turn_no        INTEGER NOT NULL REFERENCES turns (id),
    unit_id        TEXT    NOT NULL REFERENCES units (id),
    step_no        INTEGER NOT NULL,    -- order of the step within the Move
    starting_tile  INTEGER NOT NULL REFERENCES tiles (id),
    action         TEXT    NOT NULL,    -- kind of movement (Still, Follow, Scout) or direction
    ending_tile    INTEGER NOT NULL REFERENCES tiles (id),
    terrain_cd     TEXT    NOT NULL REFERENCES terrain_codes (code),
    failure_reason TEXT,                -- set only if the move failed
    parse_error    TEXT,                -- set only if the parser failed on this move
    CONSTRAINT action_check CHECK (action in ('STILL', 'SCOUT', 'N', 'NE', 'SE', 'S', 'SW', 'NW')),
    UNIQUE (clan_no, turn_no, unit_id, step_no)
);

-- --------------------------------------------------------------------------
-- Move Border Details
--
-- This table stores details about the tile borders that were found during
-- a move. The details are the border feature and the edge.
--
-- The details are always for the ending tile of the move.
CREATE TABLE move_border_details
(
    move_id   INTEGER NOT NULL REFERENCES moves (id),
    border_cd TEXT    NOT NULL REFERENCES border_codes (code),
    edge      TEXT    NOT NULL CHECK (edge in ('N', 'NE', 'SE', 'S', 'SW', 'NW')),
    PRIMARY KEY (move_id, border_cd, edge)
);

-- --------------------------------------------------------------------------
-- Move Passage Details
--
-- This table stores details about the border passages that were found during
-- a move. The details are the type of passage and the edge.
--
-- The details are always for the ending tile of the move.
CREATE TABLE move_passage_details
(
    move_id    INTEGER NOT NULL REFERENCES moves (id),
    passage_cd TEXT    NOT NULL REFERENCES passage_codes (code),
    edge       TEXT    NOT NULL CHECK (edge in ('N', 'NE', 'SE', 'S', 'SW', 'NW')),
    PRIMARY KEY (move_id, passage_cd, edge)
);

-- --------------------------------------------------------------------------
-- Move Resource Details
--
-- This table stores details about the tile resources that were found during
-- a move. The details are the type of resource and the edge.
--
-- The details are always for the ending tile of the move.
CREATE TABLE move_resource_details
(
    move_id     INTEGER NOT NULL REFERENCES moves (id),
    resource_cd TEXT    NOT NULL REFERENCES resource_codes (code),
    PRIMARY KEY (move_id, resource_cd)
);

-- --------------------------------------------------------------------------
-- Move Settlement Details
--
-- This table stores the names of settlements that were found during a move.
--
-- The details are always for the ending tile of the move.
CREATE TABLE move_settlement_details
(
    move_id INTEGER NOT NULL REFERENCES moves (id),
    name    TEXT    NOT NULL,
    PRIMARY KEY (move_id, name)
);

--  Copyright (c) 2024 Michael D Henderson. All rights reserved.

-- --------------------------------------------------------------------------
-- Move Transient Details
--
-- This table stores the units that were found during a move.
--
-- The details are always for the ending tile of the move.
CREATE TABLE move_transient_details
(
    move_id INTEGER NOT NULL REFERENCES moves (id),
    unit_id TEXT    NOT NULL REFERENCES units (id),
    PRIMARY KEY (move_id, unit_id)
);
