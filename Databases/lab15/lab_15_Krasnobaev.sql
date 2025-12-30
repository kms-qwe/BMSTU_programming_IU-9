USE master;
GO

ALTER DATABASE AirDB1 SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
GO
DROP DATABASE IF EXISTS AirDB1;
GO
CREATE DATABASE AirDB1 ON (
    NAME = AirDB1_dat,
    FILENAME = '/var/opt/mssql/data/AirDB1_dat.mdf',
    SIZE = 10MB,
    MAXSIZE = UNLIMITED,
    FILEGROWTH = 5%
)
LOG ON (
    NAME = AirDB1_log,
    FILENAME = '/var/opt/mssql/data/AirDB1_log.ldf',
    SIZE = 5MB,
    MAXSIZE = 25MB,
    FILEGROWTH = 5MB
);
GO

ALTER DATABASE AirDB2 SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
GO
DROP DATABASE IF EXISTS AirDB2;
GO
CREATE DATABASE AirDB2 ON (
    NAME = AirDB2_dat,
    FILENAME = '/var/opt/mssql/data/AirDB2_dat.mdf',
    SIZE = 10MB,
    MAXSIZE = UNLIMITED,
    FILEGROWTH = 5%
)
LOG ON (
    NAME = AirDB2_log,
    FILENAME = '/var/opt/mssql/data/AirDB2_log.ldf',
    SIZE = 5MB,
    MAXSIZE = 25MB,
    FILEGROWTH = 5MB
);
GO

USE AirDB1;
GO

DROP TRIGGER IF EXISTS dbo.trg_DeleteAirport_PreventIfRunways;
GO
DROP TABLE IF EXISTS dbo.AIRPORT;
GO

CREATE TABLE dbo.AIRPORT (
    AirportID int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    IATA char(3) NOT NULL,
    ICAO char(4) NOT NULL,
    Name nvarchar(100) NOT NULL,
    City nvarchar(100) NOT NULL,
    Timezone nvarchar(50) NOT NULL,
    CONSTRAINT UQ_AIRPORT_IATA UNIQUE (IATA),
    CONSTRAINT UQ_AIRPORT_ICAO UNIQUE (ICAO)
);
GO

CREATE TRIGGER dbo.trg_DeleteAirport_PreventIfRunways
ON dbo.AIRPORT
AFTER DELETE
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS (
        SELECT 1
        FROM deleted d
        JOIN AirDB2.dbo.RUNWAY r ON r.AirportID = d.AirportID
    )
    BEGIN
        ;THROW 51001, N'Нельзя удалить аэропорт, так как он используется в таблице RUNWAY.', 1;
    END
END;
GO

USE AirDB2;
GO

DROP TRIGGER IF EXISTS dbo.trg_DeleteRunway_PreventLast;
GO
DROP TRIGGER IF EXISTS dbo.trg_CheckRunwayAirport_Update;
GO
DROP TRIGGER IF EXISTS dbo.trg_CheckRunwayAirport_Insert;
GO
DROP TABLE IF EXISTS dbo.RUNWAY;
GO

CREATE TABLE dbo.RUNWAY (
    RunwayID int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    AirportID int NOT NULL,
    Designator char(3) NOT NULL,
    Length int NOT NULL,
    Surface nvarchar(20) NOT NULL CONSTRAINT DF_RUNWAY_Surface DEFAULT N'Asphalt',
    CONSTRAINT UQ_RUNWAY_AIRPORT_DESIGNATOR UNIQUE (AirportID, Designator)
);
GO

CREATE TRIGGER dbo.trg_CheckRunwayAirport_Insert
ON dbo.RUNWAY
AFTER INSERT
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS (
        SELECT 1
        FROM inserted i
        WHERE NOT EXISTS (
            SELECT 1
            FROM AirDB1.dbo.AIRPORT a
            WHERE a.AirportID = i.AirportID
        )
    )
    BEGIN
        ;THROW 51002, N'Указанный AirportID не существует.', 1;
    END
END;
GO

CREATE TRIGGER dbo.trg_CheckRunwayAirport_Update
ON dbo.RUNWAY
AFTER UPDATE
AS
BEGIN
    SET NOCOUNT ON;

    IF UPDATE(AirportID)
    BEGIN
        IF EXISTS (
            SELECT 1
            FROM inserted i
            WHERE NOT EXISTS (
                SELECT 1
                FROM AirDB1.dbo.AIRPORT a
                WHERE a.AirportID = i.AirportID
            )
        )
        BEGIN
            ;THROW 51003, N'Указанный AirportID не существует.', 1;
        END
    END
END;
GO

CREATE TRIGGER dbo.trg_DeleteRunway_PreventLast
ON dbo.RUNWAY
AFTER DELETE
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS (
        SELECT 1
        FROM deleted d
        JOIN AirDB1.dbo.AIRPORT ap ON ap.AirportID = d.AirportID
        WHERE NOT EXISTS (
            SELECT 1
            FROM dbo.RUNWAY r
            WHERE r.AirportID = ap.AirportID
                  AND r.RunwayID <> d.RunwayID
        )
    )
    BEGIN
        ;THROW 51004, N'Нельзя удалить последнюю ВПП аэропорта.', 1;
    END
END;
GO

USE AirDB1;
GO

INSERT INTO dbo.AIRPORT (IATA, ICAO, Name, City, Timezone)
VALUES
    ('HEL', 'EFHK', N'Helsinki-Vantaa', N'Helsinki', N'Europe/Helsinki'),
    ('VKO', 'UUWW', N'Vnukovo', N'Moscow', N'Europe/Moscow');
GO

USE AirDB2;
GO

INSERT INTO dbo.RUNWAY (AirportID, Designator, Length, Surface)
VALUES
    (1, '04L', 3500, N'Asphalt'),
    (1, '04R', 3440, N'Asphalt'),
    (2, '06', 3060, N'Concrete');
GO

SELECT * FROM AirDB1.dbo.AIRPORT;
GO
SELECT * FROM AirDB2.dbo.RUNWAY;
GO

UPDATE AirDB2.dbo.RUNWAY
SET Length = 3600
WHERE AirportID = 1 AND Designator = '04L';
GO

SELECT * FROM AirDB2.dbo.RUNWAY;
GO

UPDATE AirDB1.dbo.AIRPORT
SET Name = N'Helsinki Airport'
WHERE IATA = 'HEL';
GO

SELECT * FROM AirDB1.dbo.AIRPORT;
GO

DELETE FROM AirDB2.dbo.RUNWAY
WHERE AirportID = 1 AND Designator = '04R';
GO

SELECT * FROM AirDB2.dbo.RUNWAY;
GO
