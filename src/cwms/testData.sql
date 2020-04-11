-- positions
insert into positions (json_position) values ('{"aisle":"1a", "block":"1", "slot":"1"}');
insert into positions (json_position) values ('{"aisle":"1a", "block":"1", "slot":"2"}');
insert into positions (json_position) values ('{"aisle":"1a", "block":"1", "slot":"3"}');
insert into positions (json_position) values ('{"aisle":"1a", "block":"2", "slot":"1"}');
insert into positions (json_position) values ('{"aisle":"1a", "block":"2", "slot":"2"}');
insert into positions (json_position) values ('{"aisle":"1a", "block":"2", "slot":"3"}');

insert into positions (json_position) values ('{"aisle":"1b", "block":"1", "slot":"1"}');
insert into positions (json_position) values ('{"aisle":"1b", "block":"1", "slot":"2"}');
insert into positions (json_position) values ('{"aisle":"1b", "block":"1", "slot":"3"}');
insert into positions (json_position) values ('{"aisle":"1b", "block":"2", "slot":"1"}');
insert into positions (json_position) values ('{"aisle":"1b", "block":"2", "slot":"2"}');
insert into positions (json_position) values ('{"aisle":"1b", "block":"2", "slot":"3"}');

insert into positions (json_position) values ('{"aisle":"2a", "block":"1", "slot":"1"}');
insert into positions (json_position) values ('{"aisle":"2a", "block":"1", "slot":"2"}');
insert into positions (json_position) values ('{"aisle":"2a", "block":"1", "slot":"3"}');
insert into positions (json_position) values ('{"aisle":"2a", "block":"2", "slot":"1"}');
insert into positions (json_position) values ('{"aisle":"2a", "block":"2", "slot":"2"}');
insert into positions (json_position) values ('{"aisle":"2a", "block":"2", "slot":"3"}');
insert into positions (json_position) values ('{"aisle":"2a", "block":"3", "slot":"1"}');
insert into positions (json_position) values ('{"aisle":"2a", "block":"3", "slot":"2"}');
insert into positions (json_position) values ('{"aisle":"2a", "block":"3", "slot":"3"}');

insert into positions (json_position) values ('{"aisle":"2b", "block":"1", "slot":"1"}');
insert into positions (json_position) values ('{"aisle":"2b", "block":"1", "slot":"2"}');
insert into positions (json_position) values ('{"aisle":"2b", "block":"1", "slot":"3"}');
insert into positions (json_position) values ('{"aisle":"2b", "block":"2", "slot":"1"}');
insert into positions (json_position) values ('{"aisle":"2b", "block":"2", "slot":"2"}');
insert into positions (json_position) values ('{"aisle":"2b", "block":"2", "slot":"3"}');
insert into positions (json_position) values ('{"aisle":"2b", "block":"3", "slot":"1"}');
insert into positions (json_position) values ('{"aisle":"2b", "block":"3", "slot":"2"}');
insert into positions (json_position) values ('{"aisle":"2b", "block":"3", "slot":"3"}');

-- items
insert into items (sku, discrepancy) values ("000SKU001", "");
insert into items (sku, discrepancy) values ("empty", "");
insert into items (sku, discrepancy) values ("000SKU003", "missing");
insert into items (sku, discrepancy) values ("000SKU004", "missing");
insert into items (sku, discrepancy) values ("000SKU005", "moved");
insert into items (sku, discrepancy) values ("000SKU006", "" );
insert into items (sku, discrepancy) values ("000SKU007", "");
insert into items (sku, discrepancy) values ("empty", "");
insert into items (sku, discrepancy) values ("000SKU008", "missing");
insert into items (sku, discrepancy) values ("000SKU009", "missing");
insert into items (sku, discrepancy) values ("000SKU010", "moved");
insert into items (sku, discrepancy) values ("000SKU011", "" );

-- inventory
-- insert into item (timestamp, itemId, positionId) values ("YYYY-MM-DD HH:MM:SS.SSS", 1, 1);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.000", "2020-04-04 19:22:45.000", 1, 1);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.001", "2020-04-04 19:22:45.001", 2, 1);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.002", "2020-04-04 19:22:45.002", 3, 2);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.003", "2020-04-04 19:22:45.003", 4, 4);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.004", "2020-04-04 19:22:45.004", 5, 7);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.005", "2020-04-04 19:22:45.005", 6, 9);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.006", "2020-04-04 19:22:45.006", 7, 14);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.007", "2020-04-04 19:22:45.007", 8, 19);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.007", "2020-04-04 19:22:45.007", 9, 15);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.007", "2020-04-04 19:22:45.007", 10, 6);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.007", "2020-04-04 19:22:45.007", 11, 24);
insert into inventory (startTime, stopTime, itemId, positionId) values ("2020-04-04 19:22:45.007", "2020-04-04 19:22:45.007", 12, 23);


