USE master;
GO

IF DB_ID(N'Lab11') IS NOT NULL
BEGIN
    ALTER DATABASE Lab11 SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
    DROP DATABASE Lab11;
END;
GO

CREATE DATABASE Lab11
ON
(
    NAME = Lab11_dat,
    FILENAME = '/var/opt/mssql/data/Lab11_dat.mdf',
    SIZE = 12MB,
    MAXSIZE = UNLIMITED,
    FILEGROWTH = 4%
)
LOG ON
(
    NAME = Lab11_log,
    FILENAME = '/var/opt/mssql/data/Lab11_log.ldf',
    SIZE = 6MB,
    MAXSIZE = 30MB,
    FILEGROWTH = 6MB
);
GO

USE Lab11;
GO

DROP TABLE IF EXISTS dbo.FLIGHT;
DROP TABLE IF EXISTS dbo.AIRCRAFT;
DROP TABLE IF EXISTS dbo.RUNWAY;
DROP TABLE IF EXISTS dbo.AIRPORT;
DROP TABLE IF EXISTS dbo.AIRLINE;
GO

CREATE TABLE dbo.AIRLINE (
    AirlineID   int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    IATA        char(2)       NOT NULL,
    ICAO        char(3)       NOT NULL,
    Name        nvarchar(100) NOT NULL,
    Country     nvarchar(100) NOT NULL,
    CONSTRAINT UQ_AIRLINE_IATA UNIQUE (IATA),
    CONSTRAINT UQ_AIRLINE_ICAO UNIQUE (ICAO),
    CONSTRAINT CHK_AIRLINE_Name_NotEmpty CHECK (LEN(LTRIM(RTRIM(Name))) > 0),
    CONSTRAINT CHK_AIRLINE_Country_NotEmpty CHECK (LEN(LTRIM(RTRIM(Country))) > 0)
);
GO

CREATE TABLE dbo.AIRPORT (
    AirportID   int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    IATA        char(3)       NOT NULL,
    ICAO        char(4)       NOT NULL,
    Name        nvarchar(100) NOT NULL,
    City        nvarchar(100) NOT NULL,
    Timezone    nvarchar(50)  NOT NULL,
    CONSTRAINT UQ_AIRPORT_IATA UNIQUE (IATA),
    CONSTRAINT UQ_AIRPORT_ICAO UNIQUE (ICAO),
    CONSTRAINT CHK_AIRPORT_Name_NotEmpty CHECK (LEN(LTRIM(RTRIM(Name))) > 0),
    CONSTRAINT CHK_AIRPORT_City_NotEmpty CHECK (LEN(LTRIM(RTRIM(City))) > 0),
    CONSTRAINT CHK_AIRPORT_Timezone_NotEmpty CHECK (LEN(LTRIM(RTRIM(Timezone))) > 0)
);
GO

CREATE TABLE dbo.RUNWAY (
    RunwayID    int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    AirportID   int           NOT NULL,
    Designator  char(3)       NOT NULL,
    Length      int           NOT NULL,
    Surface     nvarchar(20)  NOT NULL CONSTRAINT DF_RUNWAY_Surface DEFAULT N'Asphalt',
    CONSTRAINT UQ_RUNWAY_AIRPORT_DESIGNATOR UNIQUE (AirportID, Designator),
    CONSTRAINT FK_RUNWAY_AIRPORT FOREIGN KEY (AirportID) REFERENCES dbo.AIRPORT (AirportID) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT CHK_RUNWAY_Length_Positive CHECK (Length > 0),
    CONSTRAINT CHK_RUNWAY_Surface_NotEmpty CHECK (LEN(LTRIM(RTRIM(Surface))) > 0)
);
GO

CREATE TABLE dbo.AIRCRAFT (
    AircraftID   int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    TailNumber   nvarchar(50)  NOT NULL,
    Model        nvarchar(100) NOT NULL,
    SeatCapacity int           NOT NULL,
    AirlineID    int           NOT NULL,
    CONSTRAINT UQ_AIRCRAFT_TAILNUMBER UNIQUE (TailNumber),
    CONSTRAINT FK_AIRCRAFT_AIRLINE FOREIGN KEY (AirlineID) REFERENCES dbo.AIRLINE (AirlineID) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT CHK_AIRCRAFT_SeatCapacity_Positive CHECK (SeatCapacity > 0),
    CONSTRAINT CHK_AIRCRAFT_Model_NotEmpty CHECK (LEN(LTRIM(RTRIM(Model))) > 0)
);
GO

CREATE TABLE dbo.FLIGHT (
    FlightID      int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    AirlineID     int           NOT NULL,
    FlightNumber  int           NOT NULL,
    OperatingDate datetime      NOT NULL,
    SchedDepTime  datetime      NOT NULL,
    SchedArrTime  datetime      NOT NULL,
    Status        nvarchar(20)  NOT NULL CONSTRAINT DF_FLIGHT_Status DEFAULT N'SCHEDULED',
    DepRunwayID   int           NOT NULL,
    ArrRunwayID   int           NOT NULL,
    AircraftID    int           NOT NULL,
    CONSTRAINT UQ_FLIGHT_AIRLINE_FLIGHTNUM_OPERDATE UNIQUE (AirlineID, FlightNumber, OperatingDate),
    CONSTRAINT FK_FLIGHT_AIRLINE FOREIGN KEY (AirlineID) REFERENCES dbo.AIRLINE (AirlineID) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_FLIGHT_DEPRUNWAY FOREIGN KEY (DepRunwayID) REFERENCES dbo.RUNWAY (RunwayID) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_FLIGHT_ARRRUNWAY FOREIGN KEY (ArrRunwayID) REFERENCES dbo.RUNWAY (RunwayID) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_FLIGHT_AIRCRAFT FOREIGN KEY (AircraftID) REFERENCES dbo.AIRCRAFT (AircraftID) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT CHK_FLIGHT_ArrAfterDep CHECK (SchedArrTime > SchedDepTime),
    CONSTRAINT CHK_FLIGHT_Status_NotEmpty CHECK (LEN(LTRIM(RTRIM(Status))) > 0)
);
GO

ALTER TABLE dbo.AIRLINE
ADD ContactEmail nvarchar(100) NULL;
GO

ALTER TABLE dbo.AIRPORT
ADD Country nvarchar(100) NOT NULL CONSTRAINT DF_AIRPORT_Country DEFAULT N'Unknown';
GO

DROP VIEW IF EXISTS dbo.vw_AirportWithRunways;
DROP VIEW IF EXISTS dbo.vw_AirlineFleet;
DROP VIEW IF EXISTS dbo.vw_FlightRoute;
GO

CREATE VIEW dbo.vw_AirportWithRunways
AS
SELECT
    a.AirportID,
    a.IATA,
    a.ICAO,
    a.Name AS AirportName,
    a.City,
    a.Timezone,
    r.RunwayID,
    r.Designator,
    r.Length,
    r.Surface
FROM dbo.AIRPORT a
LEFT JOIN dbo.RUNWAY r ON a.AirportID = r.AirportID;
GO

CREATE VIEW dbo.vw_AirlineFleet
AS
SELECT
    al.AirlineID,
    al.Name AS AirlineName,
    ac.AircraftID,
    ac.TailNumber,
    ac.Model,
    ac.SeatCapacity
FROM dbo.AIRLINE al
LEFT JOIN dbo.AIRCRAFT ac ON ac.AirlineID = al.AirlineID;
GO

CREATE VIEW dbo.vw_FlightRoute
AS
SELECT
    f.FlightID,
    a.Name AS AirlineName,
    f.FlightNumber,
    f.OperatingDate,
    da.IATA AS FromAirport,
    aa.IATA AS ToAirport,
    ac.TailNumber,
    f.Status
FROM dbo.FLIGHT f
JOIN dbo.AIRLINE a ON f.AirlineID = a.AirlineID
JOIN dbo.AIRCRAFT ac ON f.AircraftID = ac.AircraftID
JOIN dbo.RUNWAY dr ON f.DepRunwayID = dr.RunwayID
JOIN dbo.RUNWAY ar ON f.ArrRunwayID = ar.RunwayID
JOIN dbo.AIRPORT da ON dr.AirportID = da.AirportID
JOIN dbo.AIRPORT aa ON ar.AirportID = aa.AirportID;
GO


DROP TRIGGER IF EXISTS trg_InsertAircraft;
DROP TRIGGER IF EXISTS trg_DeleteAircraft;
DROP TRIGGER IF EXISTS trg_UpdateAircraftAirline;
DROP TRIGGER IF EXISTS trg_DeleteRunway_CascadeFlights;
DROP TRIGGER IF EXISTS trg_DeleteAircraft_CascadeFlights;
DROP TRIGGER IF EXISTS trg_InsertAirportWithRunways;
DROP TRIGGER IF EXISTS trg_InsertAirportWithRunways;
GO

CREATE TRIGGER trg_InsertAirportWithRunways
ON dbo.vw_AirportWithRunways
INSTEAD OF INSERT
AS
BEGIN
    SET NOCOUNT ON;

    DECLARE @old VARBINARY(128) = CONTEXT_INFO();
    DECLARE @magic VARBINARY(128) = 0x01020304;

    SET CONTEXT_INFO @magic;

    INSERT INTO dbo.AIRPORT (IATA, ICAO, Name, City, Timezone)
    SELECT DISTINCT i.IATA, i.ICAO, i.AirportName, i.City, i.Timezone
    FROM INSERTED i
    WHERE NOT EXISTS (
        SELECT 1
        FROM dbo.AIRPORT a
        WHERE a.IATA = i.IATA AND a.ICAO = i.ICAO
    );

    INSERT INTO dbo.RUNWAY (AirportID, Designator, Length, Surface)
    SELECT a.AirportID, i.Designator, i.Length, ISNULL(i.Surface, N'Asphalt')
    FROM INSERTED i
    JOIN dbo.AIRPORT a ON a.IATA = i.IATA AND a.ICAO = i.ICAO;

    SET CONTEXT_INFO @old;
END;
GO

CREATE TRIGGER trg_PreventInsertAirport
ON dbo.AIRPORT
INSTEAD OF INSERT
AS
BEGIN
    DECLARE @info VARBINARY(128);
    SELECT @info = CONTEXT_INFO();


    IF @info = 0x01020304
    BEGIN
        INSERT INTO dbo.AIRPORT (IATA, ICAO, Name, City, Timezone, Country)
        SELECT IATA, ICAO, Name, City, Timezone, Country
        FROM INSERTED;
        RETURN;
    END;

    RAISERROR(N'Нельзя добавлять аэропорт напрямую. Используйте представление vw_AirportWithRunways.', 16, 1);
    ROLLBACK TRANSACTION;
END;
GO

CREATE TRIGGER trg_DeleteRunway_CascadeFlights
ON dbo.RUNWAY
AFTER DELETE
AS
BEGIN
    SET NOCOUNT ON;

    DELETE FROM dbo.FLIGHT
    WHERE DepRunwayID IN (SELECT RunwayID FROM deleted)
       OR ArrRunwayID IN (SELECT RunwayID FROM deleted);
END;
GO

CREATE TRIGGER trg_DeleteAircraft_CascadeFlights
ON dbo.AIRCRAFT
AFTER DELETE
AS
BEGIN
    SET NOCOUNT ON;

    DELETE FROM dbo.FLIGHT
    WHERE AircraftID IN (SELECT AircraftID FROM deleted);
END;
GO

CREATE TRIGGER trg_DeleteRunway_PreventLast
ON dbo.RUNWAY
AFTER DELETE
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS (
        SELECT 1
        FROM DELETED d
        JOIN dbo.AIRPORT ap ON ap.AirportID = d.AirportID
        WHERE NOT EXISTS (
            SELECT 1
            FROM dbo.RUNWAY r
            WHERE r.AirportID = ap.AirportID
                  AND r.RunwayID <> d.RunwayID
        )
    )
    BEGIN
        RAISERROR(N'Нельзя удалить последнюю ВПП аэропорта.',16,1);
        ROLLBACK TRANSACTION;
        RETURN;
    END;
END;
GO

DROP PROCEDURE IF EXISTS dbo.GetFlightsByAirline;
GO

CREATE PROCEDURE dbo.GetFlightsByAirline
    @AirlineId int
AS
BEGIN
    SET NOCOUNT ON;

    SELECT
        f.FlightID,
        a.Name AS AirlineName,
        f.FlightNumber,
        f.OperatingDate,
        f.SchedDepTime,
        f.SchedArrTime,
        f.Status
    FROM dbo.FLIGHT f
    JOIN dbo.AIRLINE a ON f.AirlineID = a.AirlineID
    WHERE f.AirlineID = @AirlineId;
END;
GO

DROP FUNCTION IF EXISTS dbo.fn_GetFlightCountForAircraft;
DROP FUNCTION IF EXISTS dbo.fn_GetFlightsByAirport;
GO

CREATE FUNCTION dbo.fn_GetFlightCountForAircraft
(
    @AircraftId int
)
RETURNS int
AS
BEGIN
    DECLARE @cnt int;
    SELECT @cnt = COUNT(*) FROM dbo.FLIGHT WHERE AircraftID = @AircraftId;
    RETURN @cnt;
END;
GO

CREATE FUNCTION dbo.fn_GetFlightsByAirport
(
    @AirportId int
)
RETURNS TABLE
AS
RETURN
(
    SELECT
        f.FlightID,
        f.FlightNumber,
        f.OperatingDate,
        f.Status,
        da.IATA AS DepAirport,
        aa.IATA AS ArrAirport
    FROM dbo.FLIGHT f
    JOIN dbo.RUNWAY dr ON f.DepRunwayID = dr.RunwayID
    JOIN dbo.RUNWAY ar ON f.ArrRunwayID = ar.RunwayID
    JOIN dbo.AIRPORT da ON dr.AirportID = da.AirportID
    JOIN dbo.AIRPORT aa ON ar.AirportID = aa.AirportID
    WHERE da.AirportID = @AirportId OR aa.AirportID = @AirportId
);
GO

INSERT INTO dbo.AIRLINE (IATA, ICAO, Name, Country, ContactEmail)
VALUES
('SU','AFL',N'Aeroflot',N'Russia',N'info@aeroflot.ru'),
('UT','UTA',N'UTair',N'Russia',N'contact@utair.ru'),
('DP','PBD',N'Pobeda',N'Russia',NULL);
GO

INSERT INTO dbo.vw_AirportWithRunways
    (AirportID, IATA, ICAO, AirportName, City, Timezone, RunwayID, Designator, Length, Surface)
VALUES
    (NULL,'DME','UUDD',N'Domodedovo',N'Moscow',N'UTC+3',NULL,'14L',3800,N'Concrete');
GO

INSERT INTO dbo.RUNWAY (AirportID, Designator, Length, Surface)
VALUES
(1,'06L',3500,N'Concrete'),
(1,'06R',3700,N'Asphalt');
GO

INSERT INTO dbo.AIRCRAFT (TailNumber, Model, SeatCapacity, AirlineID)
VALUES
(N'RA-89001',N'SSJ-100',100,1),
(N'RA-73801',N'B737-800',160,1),
(N'VP-BLA',N'A320',150,2),
(N'VP-BBB',N'B737-800',180,3);
GO

INSERT INTO dbo.FLIGHT (AirlineID, FlightNumber, OperatingDate, SchedDepTime, SchedArrTime, Status, DepRunwayID, ArrRunwayID, AircraftID)
VALUES
(1,1001,'2024-12-01T00:00:00','2024-12-01T09:00:00','2024-12-01T10:30:00',N'SCHEDULED',1,3,1),
(1,1002,'2024-12-01T00:00:00','2024-12-01T13:00:00','2024-12-01T14:30:00',N'ARRIVED',2,3,2);
GO

INSERT INTO dbo.FLIGHT (AirlineID, FlightNumber, OperatingDate, SchedDepTime, SchedArrTime, Status, DepRunwayID, ArrRunwayID, AircraftID)
SELECT AirlineID,
       FlightNumber + 10,
       DATEADD(day,1,OperatingDate),
       DATEADD(hour,1,SchedDepTime),
       DATEADD(hour,1,SchedArrTime),
       N'SCHEDULED',
       DepRunwayID,
       ArrRunwayID,
       AircraftID
FROM dbo.FLIGHT
WHERE FlightID = 1;
GO

UPDATE dbo.AIRCRAFT
SET AirlineID = 2
WHERE AircraftID = 1;
GO

EXEC dbo.GetFlightsByAirline @AirlineId = 1;
GO

SELECT dbo.fn_GetFlightCountForAircraft(1) AS FlightsForAircraft;
GO

SELECT * FROM dbo.fn_GetFlightsByAirport(1);
GO

SELECT DISTINCT Country
FROM dbo.AIRLINE;
GO

SELECT Name, Country
FROM dbo.AIRLINE
ORDER BY Name ASC;
GO

SELECT TailNumber, SeatCapacity
FROM dbo.AIRCRAFT
ORDER BY SeatCapacity DESC;
GO

SELECT f.FlightNumber, a.Name AS AirlineName
FROM dbo.FLIGHT f
JOIN dbo.AIRLINE a ON f.AirlineID = a.AirlineID;
GO

SELECT al.Name AS AirlineName, ac.TailNumber
FROM dbo.AIRLINE al
LEFT JOIN dbo.AIRCRAFT ac ON ac.AirlineID = al.AirlineID;
GO

SELECT r.Designator, ap.IATA
FROM dbo.RUNWAY r
RIGHT JOIN dbo.AIRPORT ap ON r.AirportID = ap.AirportID;
GO

SELECT ap.IATA, r.Designator
FROM dbo.AIRPORT ap
FULL OUTER JOIN dbo.RUNWAY r ON ap.AirportID = r.AirportID;
GO

SELECT Name, Country
FROM dbo.AIRLINE
WHERE Name LIKE N'A%';
GO

SELECT FlightNumber, OperatingDate
FROM dbo.FLIGHT
WHERE OperatingDate BETWEEN '2024-12-01' AND '2024-12-02';
GO

SELECT FlightNumber, Status
FROM dbo.FLIGHT
WHERE Status IN (N'SCHEDULED',N'DELAYED');
GO

SELECT ap.Name, ap.IATA
FROM dbo.AIRPORT ap
WHERE EXISTS (
    SELECT 1
    FROM dbo.RUNWAY r
    JOIN dbo.FLIGHT f ON f.DepRunwayID = r.RunwayID
    WHERE r.AirportID = ap.AirportID
);
GO

SELECT Name, ContactEmail
FROM dbo.AIRLINE
WHERE ContactEmail IS NULL;
GO

SELECT a.Name AS AirlineName, COUNT(f.FlightID) AS FlightCount
FROM dbo.AIRLINE a
LEFT JOIN dbo.FLIGHT f ON f.AirlineID = a.AirlineID
GROUP BY a.Name;
GO

SELECT a.Name AS AirlineName, COUNT(f.FlightID) AS FlightCount
FROM dbo.AIRLINE a
JOIN dbo.FLIGHT f ON f.AirlineID = a.AirlineID
GROUP BY a.Name
HAVING COUNT(f.FlightID) > 1;
GO

SELECT
    a.Name AS AirlineName,
    SUM(ac.SeatCapacity) AS TotalSeats,
    MIN(ac.SeatCapacity) AS MinSeats,
    MAX(ac.SeatCapacity) AS MaxSeats,
    AVG(CAST(ac.SeatCapacity AS float)) AS AvgSeats
FROM dbo.AIRLINE a
JOIN dbo.AIRCRAFT ac ON ac.AirlineID = a.AirlineID
GROUP BY a.Name;
GO

SELECT Country AS Name
FROM dbo.AIRPORT
UNION
SELECT Country AS Name
FROM dbo.AIRLINE;
GO

SELECT Country AS Name
FROM dbo.AIRPORT
UNION ALL
SELECT Country AS Name
FROM dbo.AIRLINE;
GO

SELECT Country
FROM dbo.AIRPORT
EXCEPT
SELECT Country
FROM dbo.AIRLINE;
GO

SELECT Country
FROM dbo.AIRPORT
INTERSECT
SELECT Country
FROM dbo.AIRLINE;
GO

UPDATE dbo.AIRCRAFT
SET SeatCapacity = SeatCapacity + 10
WHERE AircraftID IN (
    SELECT AircraftID
    FROM dbo.FLIGHT
    GROUP BY AircraftID
    HAVING COUNT(*) >= 2
);
GO

DELETE FROM dbo.AIRCRAFT
WHERE NOT EXISTS (
    SELECT 1
    FROM dbo.FLIGHT f
    WHERE f.AircraftID = dbo.AIRCRAFT.AircraftID
);
GO

DELETE FROM dbo.AIRPORT
WHERE NOT EXISTS (
    SELECT 1
    FROM dbo.RUNWAY r
    WHERE r.AirportID = dbo.AIRPORT.AirportID
);
GO

EXEC dbo.GetFlightsByAirline @AirlineId = 2;
GO

SELECT * FROM dbo.vw_FlightRoute;
GO

SELECT * FROM dbo.vw_AirlineFleet;
GO
