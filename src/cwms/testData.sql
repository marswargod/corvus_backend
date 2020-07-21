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

insert into regions (name, frequency) values ("region1", 3);
insert into regions (name, frequency) values ("region2", 1);
insert into regions (name, frequency) values ("region3", 1);
insert into regions (name, frequency) values ("region4", 1);
insert into regions (name, frequency) values ("region5", 1);
insert into regions (name, frequency) values ("region6", 1);
insert into regions (name, frequency) values ("region7", 1);
insert into regions (name, frequency) values ("region8", 1);
insert into regions (name, frequency) values ("region9", 1);
insert into regions (name, frequency) values ("region10", 1);
insert into regions (name, frequency) values ("region11", 1);
insert into regions (name, frequency) values ("region12", 1);

insert into regionPositions (regionId, positionId) values (1,1);
insert into regionPositions (regionId, positionId) values (1,2);
insert into regionPositions (regionId, positionId) values (1,3);
insert into regionPositions (regionId, positionId) values (1,4);
insert into regionPositions (regionId, positionId) values (1,5);
insert into regionPositions (regionId, positionId) values (1,6);

insert into regionPositions (regionId, positionId) values (2,7);
insert into regionPositions (regionId, positionId) values (2,8);
insert into regionPositions (regionId, positionId) values (2,9);
insert into regionPositions (regionId, positionId) values (2,10);
insert into regionPositions (regionId, positionId) values (2,11);
insert into regionPositions (regionId, positionId) values (2,12);

insert into regionPositions (regionId, positionId) values (3,13);
insert into regionPositions (regionId, positionId) values (3,14);
insert into regionPositions (regionId, positionId) values (3,15);
insert into regionPositions (regionId, positionId) values (3,16);
insert into regionPositions (regionId, positionId) values (3,17);
insert into regionPositions (regionId, positionId) values (3,18);
insert into regionPositions (regionId, positionId) values (3,19);
insert into regionPositions (regionId, positionId) values (3,20);
insert into regionPositions (regionId, positionId) values (3,21);

insert into regionPositions (regionId, positionId) values (4,22);
insert into regionPositions (regionId, positionId) values (4,23);
insert into regionPositions (regionId, positionId) values (4,24);
insert into regionPositions (regionId, positionId) values (4,25);
insert into regionPositions (regionId, positionId) values (4,26);
insert into regionPositions (regionId, positionId) values (4,27);
insert into regionPositions (regionId, positionId) values (4,28);
insert into regionPositions (regionId, positionId) values (4,29);
insert into regionPositions (regionId, positionId) values (4,30);

insert into events (name, entry, regionId) values ("event1",  1, 1);
insert into events (name, entry, regionId) values ("event2",  2, 2);
insert into events (name, entry, regionId) values ("event3",  3, 3);
insert into events (name, entry, regionId) values ("event4",  4, 4);
insert into events (name, entry, regionId) values ("event5",  5, 1);
insert into events (name, entry, regionId) values ("event6",  6, 5);
insert into events (name, entry, regionId) values ("event7",  7, 6);
insert into events (name, entry, regionId) values ("event8",  8, 7);
insert into events (name, entry, regionId) values ("event9",  9, 8);
insert into events (name, entry, regionId) values ("event10", 10, 1);
insert into events (name, entry, regionId) values ("event11", 11, 9);
insert into events (name, entry, regionId) values ("event12", 12, 10);
insert into events (name, entry, regionId) values ("event13", 13, 11);
insert into events (name, entry, regionId) values ("event14", 14, 12);

insert into restrictions (name, startTime, stopTime, periodicity, regionId) values ("causeway", "10:00", "13:00", "daily", 1);
insert into restrictions (name, startTime, stopTime, periodicity, regionId) values ("crossroads", "10:00", "13:00", "daily", 2);
insert into restrictions (name, startTime, stopTime, periodicity, regionId) values ("footpath", "10:00", "13:00", "daily", 3);
insert into restrictions (name, startTime, stopTime, periodicity, regionId) values ("area57", "10:00", "13:00", "daily", 4);
insert into restrictions (name, startTime, stopTime, periodicity, regionId) values ("breezeway", "10:00", "13:00", "daily", 5);

insert into flights (time) values ("10:00");
insert into flights (time) values ("10:01");
insert into flights (time) values ("10:02");

insert into flightPositions (flightId, positionId, sku, occupancy) values (1, 1, "000SKU005", "12.1");
insert into flightPositions (flightId, positionId, sku, occupancy) values (1, 2, "000SKU006", "12.2");
insert into flightPositions (flightId, positionId, sku, occupancy) values (1, 3, "000SKU007", "12.3");
insert into flightPositions (flightId, positionId, sku, occupancy) values (2, 4, "000SKU008", "13.1");
insert into flightPositions (flightId, positionId, sku, occupancy) values (2, 5, "000SKU009", "13.2");
insert into flightPositions (flightId, positionId, sku, occupancy) values (2, 6, "000SKU010", "13.3");