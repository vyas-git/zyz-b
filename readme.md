# Minimum Purchase License Count

### Run

`go run main.go sample-small.csv`

### Run Concurrency

`go run concurency/concurrency.go sample-small.csv`

### Test

`go test`

### SQL Query

- Import csv file data into local postgress db and run below query to cross verify count for larger data sets

`SELECT
COALESCE(
SUM(
CASE
WHEN "LaptopCount" > 0 THEN GREATEST("LaptopCount", "DesktopCount" - "LaptopCount")
ELSE "DesktopCount"
END
), 0
) AS "MinCopiesRequired"
FROM (
SELECT
"UserID",
COUNT(DISTINCT CASE WHEN LOWER("ComputerType") = 'desktop' THEN "ComputerID" END) AS "DesktopCount",
COUNT(DISTINCT CASE WHEN LOWER("ComputerType") = 'laptop' THEN "ComputerID" END) AS "LaptopCount"
FROM
"schema"."tablename"
WHERE
"ApplicationID" = 374
GROUP BY
"UserID"
) AS "UserCounts";`

- Replace shema and tablename in above query
