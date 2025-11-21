USE master;
GO

IF DB_ID('Lab7') IS NOT NULL
BEGIN
    ALTER DATABASE Lab7 SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
    DROP DATABASE Lab7;
END
GO

CREATE DATABASE Lab7
ON ( NAME = Lab7_dat,
     FILENAME = '/var/opt/mssql/data/Lab7_dat.mdf',
     SIZE = 10MB,
     MAXSIZE = UNLIMITED,
     FILEGROWTH = 5%
)
LOG ON ( NAME = Lab7_log,
         FILENAME = '/var/opt/mssql/data/Lab7_log.ldf',
         SIZE = 5MB,
         MAXSIZE = 25MB,
         FILEGROWTH = 5MB
);
GO

USE Lab7;
GO

DROP TABLE IF EXISTS dbo.AIRLINE;
GO

CREATE TABLE dbo.AIRLINE (
    AirlineID   INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    IATA        CHAR(2)       NOT NULL,
    ICAO        CHAR(3)       NOT NULL,
    Name        NVARCHAR(100) NOT NULL,
    Country     NVARCHAR(100) NOT NULL,
    CONSTRAINT AK_AIRLINE_IATA UNIQUE (IATA),
    CONSTRAINT AK_AIRLINE_ICAO UNIQUE (ICAO)
);
GO

DROP TABLE IF EXISTS dbo.AIRCRAFT;
GO

CREATE TABLE dbo.AIRCRAFT (
    AircraftID   INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    TailNumber   NVARCHAR(50)  NOT NULL,
    Model        NVARCHAR(100) NOT NULL,
    SeatCapacity INT           NOT NULL,
    AirlineID    INT           NOT NULL,
    CONSTRAINT AK_AIRCRAFT_TAILNUMBER UNIQUE (TailNumber),
    CONSTRAINT FK_AIRCRAFT_AIRLINE
        FOREIGN KEY (AirlineID)
        REFERENCES dbo.AIRLINE (AirlineID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
GO

INSERT INTO dbo.AIRLINE (IATA, ICAO, Name, Country)
VALUES
    ('SU', 'AFL', N'Aeroflot',      N'Russia'),
    ('LH', 'DLH', N'Lufthansa',     N'Germany'),
    ('AF', 'AFR', N'Air France',    N'France');
GO

INSERT INTO dbo.AIRCRAFT (TailNumber, Model, SeatCapacity, AirlineID)
VALUES
    ('RA-89001', N'Sukhoi Superjet 100', 100, 1),
    ('RA-89002', N'Sukhoi Superjet 100', 108, 1),
    ('D-AIXX',   N'Airbus A320',         180, 2),
    ('F-GKXA',   N'Airbus A321',         220, 3);
GO

CREATE VIEW dbo.vw_AircraftLarge
AS
    SELECT
        AircraftID,
        TailNumber,
        Model,
        SeatCapacity,
        AirlineID
    FROM dbo.AIRCRAFT
    WHERE SeatCapacity > 150
    WITH CHECK OPTION;
GO

CREATE VIEW dbo.vw_AircraftWithAirline
AS
    SELECT
        ac.AircraftID,
        ac.TailNumber,
        ac.Model,
        ac.SeatCapacity,
        al.AirlineID,
        al.Name    AS AirlineName,
        al.IATA    AS AirlineIATA,
        al.Country AS AirlineCountry
    FROM dbo.AIRCRAFT ac
    JOIN dbo.AIRLINE  al
        ON ac.AirlineID = al.AirlineID;
GO

CREATE NONCLUSTERED INDEX IDX_Aircraft_Model
ON dbo.AIRCRAFT (Model)
INCLUDE (SeatCapacity, AirlineID, TailNumber);
GO

CREATE VIEW dbo.vw_IndexedAirline
WITH SCHEMABINDING
AS
    SELECT
        AirlineID,
        IATA,
        ICAO,
        Name,
        Country
    FROM dbo.AIRLINE;
GO

CREATE UNIQUE CLUSTERED INDEX IDX_vw_IndexedAirline
ON dbo.vw_IndexedAirline (AirlineID);
GO

SELECT * FROM dbo.vw_AircraftLarge;
SELECT * FROM dbo.vw_AircraftWithAirline;

SELECT AircraftID, Model, SeatCapacity, AirlineID
FROM dbo.AIRCRAFT
WHERE Model = N'Airbus A320';

SELECT * FROM dbo.vw_IndexedAirline;
GO
