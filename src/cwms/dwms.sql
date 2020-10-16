PRAGMA foreign_keys = ON;


DROP VIEW IF EXISTS v_schedule;
DROP VIEW IF EXISTS v_regionPosition;
DROP VIEW IF EXISTS v_aisleStats;
DROP VIEW IF EXISTS v_inventory;
DROP VIEW IF EXISTS v_restrictions;
DROP VIEW IF EXISTS v_flightList;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS restrictions;
DROP TABLE IF EXISTS positions;
DROP TABLE IF EXISTS regionPositions;
DROP TABLE IF EXISTS flightPositions;
DROP TABLE IF EXISTS regions;
DROP TABLE IF EXISTS flights;
DROP TABLE IF EXISTS items;


CREATE TABLE IF NOT EXISTS positions (
  positionId INTEGER PRIMARY KEY AUTOINCREMENT,
  json_position TEXT
);
-- Positions are stored in json e.g. {"Aisle":"1a","Shelf":"1","Slot":"1"}
-- https://www.sqlite.org/json1.html#jex
-- https://community.esri.com/groups/appstudio/blog/2018/08/21/working-with-json-in-sqlite-databases
DROP INDEX IF EXISTS idx_aisle;
CREATE INDEX idx_aisle ON positions (json_extract(json_position, '$.aisle'));

CREATE TABLE IF NOT EXISTS items (
  itemId INTEGER PRIMARY KEY AUTOINCREMENT,
  sku TEXT,
  discrepancy TEXT
);
DROP INDEX IF EXISTS idx_sku;
CREATE INDEX idx_sku ON items (sku);
DROP INDEX IF EXISTS idx_discrepancy;
CREATE INDEX idx_discrepancy ON items (discrepancy);

CREATE TABLE IF NOT EXISTS images (
  imageId INTEGER PRIMARY KEY AUTOINCREMENT,
  imageUrl TEXT
);

CREATE TABLE IF NOT EXISTS inventory (
  inventoryId INTEGER PRIMARY KEY AUTOINCREMENT,
  startTime DATETIME,
  stopTime DATETIME,
  itemId INTEGER REFERENCES items(itemId),
  positionId INTEGER REFERENCES positions(positionId),
  imageId INTEGER REFERENCES images(imageId)
);
-- Timestamps are stored using unix timestamps
-- number of seconds that have passed since midnight on the 1st January 1970, UTC time
-- https://www.sqlite.org/lang_datefunc.html
-- https://www.sqlite.org/draft/datatype3.html

CREATE VIEW v_inventory
  AS SELECT
    inventoryId,
    startTime,
    stopTime,
    items.sku AS sku,
    json_extract(positions.json_position, "$.aisle") AS aisle,
    json_extract(positions.json_position, "$.block") AS block,
    json_extract(positions.json_position, "$.slot") AS slot,
    json_extract(positions.json_position, "$.shelf") AS shelf,
    json_extract(positions.json_position, "$.displayname") AS displayName,
    items.discrepancy AS discrepancy,
    images.imageUrl AS imageUrl
  FROM
    inventory
    LEFT JOIN positions USING(positionId)
    LEFT JOIN items USING(itemId)
    LEFT JOIN images USING(imageId);

CREATE VIEW IF NOT EXISTS v_aisleStats
  AS SELECT
    aisle,
    sum(case when discrepancy != "" then 1 else 0 end) as numberException,
    sum(case when sku = "empty" then 1 else 0 end) as numberEmpty,
    sum(case when sku != "empty" and sku is not null then 1 else 0 end) as numberOccupied,
    sum(case when sku is null then 1 else 0 end) as numberUnscanned,
    max(stopTime) as lastScanned -- TODO: might not be right
  FROM
    v_inventory
  GROUP BY
    aisle;

CREATE TABLE IF NOT EXISTS regions (
  regionId INTEGER PRIMARY KEY AUTOINCREMENT,
  name string,
  frequency int
);

CREATE TABLE IF NOT EXISTS regionPositions (
  rpId INTEGER PRIMARY KEY AUTOINCREMENT,
  name string,
  regionId INTEGER REFERENCES regions(regionId),
  positionId INTEGER REFERENCES positions(positionId)
);

CREATE VIEW IF NOT EXISTS v_regionPosition
  AS
  SELECT
    regionId AS regionId,
    json_extract(positions.json_position, "$.aisle") AS aisle
  FROM
    regions
    LEFT JOIN regionPositions USING(regionId)
    LEFT JOIN positions USING(positionId);


CREATE TABLE IF NOT EXISTS events (
  eventId INTEGER PRIMARY KEY AUTOINCREMENT,
  name string,
  queue string,
  entry int,
  regionId INTEGER REFERENCES regions(regionId)
);

CREATE TABLE IF NOT EXISTS restrictions (
  restrictionId INTEGER PRIMARY KEY AUTOINCREMENT,
  name string,
  startDate DATETIME,
  stopDate  DATETIME,
  startTime DATETIME,
  stopTime  DATETIME,
  periodicityNum int,
  periodicity string,
  regionId INTEGER REFERENCES regions(regionId)
  -- CHECK (periodicity IN ('weekdays','weekends','everyday','monday','tuesday','wednesday','thursday','friday','saturday','sunday'))
);

CREATE VIEW IF NOT EXISTS v_restrictions
  AS SELECT
    restrictionId AS restrictionId,
    startDate DATETIME,
    stopDate  DATETIME,
    json_extract(positions.json_position, "$.aisle") AS aisle
  FROM
    restrictions
    LEFT JOIN regions USING(regionId)
    LEFT JOIN regionPositions USING(regionId)
    LEFT JOIN positions USING(positionId);

CREATE VIEW IF NOT EXISTS v_schedule
  AS SELECT
    entry AS entry,
    queue AS queue,
    regions.name AS region,
    regions.frequency AS frequency,
    restrictions.name AS restriction,
    restrictions.startTime AS startTime,
    restrictions.stopTime AS stopTime,
    restrictions.periodicity AS periodicity,
    json_extract(positions.json_position, "$.aisle") AS aisle,
    json_extract(positions.json_position, "$.block") AS block,
    json_extract(positions.json_position, "$.slot") AS slot
  FROM
    events
    LEFT JOIN regions USING(regionId)
    LEFT JOIN regionPositions USING(regionId)
    LEFT JOIN restrictions USING(regionId)
    LEFT JOIN positions USING(positionId);

CREATE TABLE IF NOT EXISTS flights (
  flightId  INTEGER PRIMARY KEY AUTOINCREMENT,
  time DATETIME
);

CREATE TABLE IF NOT EXISTS flightPositions (
  fpId  INTEGER PRIMARY KEY AUTOINCREMENT,
  sku text,
  occupancy text,
  flightId INTEGER REFERENCES flights(flightId),
  positionId INTEGER REFERENCES positions(positionId)
);

CREATE VIEW v_flightList
  AS SELECT
    flightId AS flightId,
    time(time) AS time,
    flightPositions.sku AS sku,
    flightPositions.occupancy AS occupancy,
    json_extract(positions.json_position, "$.aisle") AS aisle,
    json_extract(positions.json_position, "$.block") AS block,
    json_extract(positions.json_position, "$.slot") AS slot
  FROM
    flights
    LEFT JOIN flightPositions USING(flightId)
    LEFT JOIN positions USING(positionId);
