PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS positions;
CREATE TABLE IF NOT EXISTS positions (
  positionId INTEGER PRIMARY KEY AUTOINCREMENT,
  json_position TEXT
);
-- Positions are stored in json e.g. {"Aisle":"1a","Shelf":"1","Slot":"1"}
-- https://www.sqlite.org/json1.html#jex
-- https://community.esri.com/groups/appstudio/blog/2018/08/21/working-with-json-in-sqlite-databases
CREATE INDEX idx_aisle ON positions (json_extract(json_position, '$.aisle'));

DROP TABLE IF EXISTS items;
CREATE TABLE IF NOT EXISTS items (
  itemId INTEGER PRIMARY KEY AUTOINCREMENT,
  sku TEXT,
  discrepancy TEXT
);
CREATE INDEX idx_sku ON items (sku);
CREATE INDEX idx_discrepany ON items (discrepancy);

DROP TABLE IF EXISTS inventory;
CREATE TABLE IF NOT EXISTS inventory (
  inventoryId INTEGER PRIMARY KEY AUTOINCREMENT,
  startTime DATETIME, 
  stopTime DATETIME,
  itemId INTEGER REFERENCES items(itemId),
  positionId INTEGER REFERENCES positions(positionId)
);
-- Timestamps are stored using unix timestamps
-- number of seconds that have passed since midnight on the 1st January 1970, UTC time
-- https://www.sqlite.org/lang_datefunc.html
-- https://www.sqlite.org/draft/datatype3.html

CREATE VIEW v_inventory
AS
SELECT
  inventoryId,
  startTime,
  stopTime,
  items.sku AS sku,
  json_extract(positions.json_position, "$.aisle") AS aisle,
  json_extract(positions.json_position, "$.block") AS block,
  json_extract(positions.json_position, "$.slot") AS slot,
  items.discrepancy AS discrepancy
FROM
  inventory
  LEFT JOIN positions USING(positionId)
  LEFT JOIN items USING(itemId);