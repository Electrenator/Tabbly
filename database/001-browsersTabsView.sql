BEGIN TRANSACTION;

-- Uses no indexes given that increases the DB size quite substantially. This won't be
-- executed often or with time sensitivity anyways. Mostly only for manually viewing
-- the data and maybe a future, currently unplanned, feature.
CREATE VIEW BrowserStatsView AS
SELECT
    Browser.name,
    Entry.timestamp,
    COUNT(Window.openTabs) AS totalWindows,
    SUM(Window.openTabs) as totalTabs
FROM `Entry` INNER JOIN (
	SELECT * FROM Browser
) AS Browser ON Browser.id = Entry.browserId LEFT JOIN (
	SELECT * FROM Window
) AS Window ON Window.entryId = Entry.id GROUP BY Entry.timestamp, Browser.name;

UPDATE `Database` SET `version` = 1;

COMMIT;
