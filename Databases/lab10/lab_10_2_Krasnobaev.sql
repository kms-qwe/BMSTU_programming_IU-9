USE master;
GO

IF DB_ID(N'Lab10') IS NOT NULL
BEGIN
    ALTER DATABASE Lab10 SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
    DROP DATABASE Lab10;
END;
GO

CREATE DATABASE Lab10
ON (
    NAME = Lab10_dat,
    FILENAME = '/var/opt/mssql/data/Lab10_dat.mdf',
    SIZE = 10MB,
    MAXSIZE = UNLIMITED,
    FILEGROWTH = 5%
)
LOG ON (
    NAME = Lab10_log,
    FILENAME = '/var/opt/mssql/data/Lab10_log.ldf',
    SIZE = 5MB,
    MAXSIZE = 5GB,
    FILEGROWTH = 5MB
);
GO

USE Lab10;
GO

DROP TABLE IF EXISTS dbo.AIRCRAFT;
DROP TABLE IF EXISTS dbo.AIRLINE;
GO

CREATE TABLE dbo.AIRLINE (
    AirlineID   INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    IATA        CHAR(2)       NOT NULL,
    ICAO        CHAR(3)       NOT NULL,
    Name        NVARCHAR(100) NOT NULL,
    Country     NVARCHAR(100) NOT NULL,
    CONSTRAINT UQ_AIRLINE_IATA UNIQUE (IATA),
    CONSTRAINT UQ_AIRLINE_ICAO UNIQUE (ICAO)
);
GO

CREATE TABLE dbo.AIRCRAFT (
    AircraftID   INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    TailNumber   NVARCHAR(50)  NOT NULL,
    Model        NVARCHAR(100) NOT NULL,
    SeatCapacity INT           NOT NULL,
    AirlineID    INT           NOT NULL,
    CONSTRAINT UQ_AIRCRAFT_TailNumber UNIQUE (TailNumber),
    CONSTRAINT FK_AIRCRAFT_AIRLINE FOREIGN KEY (AirlineID)
        REFERENCES dbo.AIRLINE (AirlineID) ON DELETE CASCADE ON UPDATE CASCADE
);
GO

INSERT INTO dbo.AIRLINE (IATA, ICAO, Name, Country)
VALUES
    ('SU', 'AFL', N'Aeroflot',        N'Russia'),
    ('S7', 'SBI', N'S7 Airlines',     N'Russia'),
    ('LH', 'DLH', N'Lufthansa',       N'Germany'),
    ('AF', 'AFR', N'Air France',      N'France'),
    ('BA', 'BAW', N'British Airways', N'United Kingdom');
GO

INSERT INTO dbo.AIRCRAFT (TailNumber, Model, SeatCapacity, AirlineID)
VALUES
    (N'RA-001', N'SuperJet 100', 100, 1),
    (N'RA-002', N'A320',         180, 1),
    (N'VP-B01', N'A321',         200, 2),
    (N'D-AI01', N'A320',         180, 3);
GO

USE Lab10;
GO

BEGIN TRANSACTION
UPDATE dbo.AIRCRAFT
SET SeatCapacity = 130
WHERE TailNumber = N'RA-001';
WAITFOR DELAY '00:00:07';
SELECT * FROM dbo.AIRCRAFT;
SELECT * FROM sys.dm_tran_locks
WHERE resource_database_id = DB_ID(N'Lab10');
ROLLBACK TRANSACTION;
GO

BEGIN TRANSACTION
UPDATE dbo.AIRCRAFT
SET SeatCapacity = 190
WHERE TailNumber = N'RA-002';
WAITFOR DELAY '00:00:07';
SELECT * FROM dbo.AIRCRAFT;
SELECT * FROM sys.dm_tran_locks
WHERE resource_database_id = DB_ID(N'Lab10');
COMMIT TRANSACTION;
GO

BEGIN TRANSACTION
INSERT INTO dbo.AIRCRAFT (TailNumber, Model, SeatCapacity, AirlineID)
VALUES (N'RA-009', N'A321', 220, 1);
WAITFOR DELAY '00:00:07';
SELECT * FROM dbo.AIRCRAFT;
SELECT * FROM sys.dm_tran_locks
WHERE resource_database_id = DB_ID(N'Lab10');
COMMIT TRANSACTION;
GO

BEGIN TRANSACTION
INSERT INTO dbo.AIRCRAFT (TailNumber, Model, SeatCapacity, AirlineID)
VALUES (N'RA-010', N'A321', 230, 1);
WAITFOR DELAY '00:00:07';
SELECT * FROM dbo.AIRCRAFT;
SELECT * FROM sys.dm_tran_locks
WHERE resource_database_id = DB_ID(N'Lab10');
COMMIT TRANSACTION;
GO
