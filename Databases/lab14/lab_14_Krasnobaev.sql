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

DROP TABLE IF EXISTS dbo.AIRLINE_P1;
GO
CREATE TABLE dbo.AIRLINE_P1 (
    AirlineID int NOT NULL,
    IATA char(2) NOT NULL,
    ICAO char(3) NOT NULL,
    CONSTRAINT PK_AIRLINE_P1 PRIMARY KEY (AirlineID),
    CONSTRAINT UQ_AIRLINE_P1_IATA UNIQUE (IATA),
    CONSTRAINT UQ_AIRLINE_P1_ICAO UNIQUE (ICAO)
);
GO

USE AirDB2;
GO

DROP TABLE IF EXISTS dbo.AIRLINE_P2;
GO
CREATE TABLE dbo.AIRLINE_P2 (
    AirlineID int NOT NULL,
    Name nvarchar(100) NOT NULL,
    Country nvarchar(100) NOT NULL,
    CONSTRAINT PK_AIRLINE_P2 PRIMARY KEY (AirlineID)
);
GO

USE AirDB1;
GO

DROP VIEW IF EXISTS dbo.AIRLINE_View;
GO
CREATE VIEW dbo.AIRLINE_View
AS
    SELECT
        p1.AirlineID,
        p1.IATA,
        p1.ICAO,
        p2.Name,
        p2.Country
    FROM AirDB1.dbo.AIRLINE_P1 p1
        INNER JOIN AirDB2.dbo.AIRLINE_P2 p2 ON p1.AirlineID = p2.AirlineID;
GO

DROP TRIGGER IF EXISTS dbo.trg_UpdateAIRLINE_View;
GO
CREATE TRIGGER dbo.trg_UpdateAIRLINE_View
ON dbo.AIRLINE_View
INSTEAD OF UPDATE
AS
BEGIN
    SET NOCOUNT ON;

    IF UPDATE(AirlineID)
    BEGIN
        ;THROW 50001, N'Нельзя изменить столбец AirlineID', 1;
    END

    UPDATE p1
    SET
        p1.IATA = inserted.IATA,
        p1.ICAO = inserted.ICAO
    FROM AirDB1.dbo.AIRLINE_P1 p1
        INNER JOIN inserted ON p1.AirlineID = inserted.AirlineID;

    UPDATE p2
    SET
        p2.Name = inserted.Name,
        p2.Country = inserted.Country
    FROM AirDB2.dbo.AIRLINE_P2 p2
        INNER JOIN inserted ON p2.AirlineID = inserted.AirlineID;
END;
GO

DROP TRIGGER IF EXISTS dbo.trg_InsertAIRLINE_View;
GO
CREATE TRIGGER dbo.trg_InsertAIRLINE_View
ON dbo.AIRLINE_View
INSTEAD OF INSERT
AS
BEGIN
    SET NOCOUNT ON;

    INSERT INTO AirDB1.dbo.AIRLINE_P1 (AirlineID, IATA, ICAO)
    SELECT AirlineID, IATA, ICAO
    FROM inserted;

    INSERT INTO AirDB2.dbo.AIRLINE_P2 (AirlineID, Name, Country)
    SELECT AirlineID, Name, Country
    FROM inserted;
END;
GO

DROP TRIGGER IF EXISTS dbo.trg_DeleteAIRLINE_View;
GO
CREATE TRIGGER dbo.trg_DeleteAIRLINE_View
ON dbo.AIRLINE_View
INSTEAD OF DELETE
AS
BEGIN
    SET NOCOUNT ON;

    DELETE FROM AirDB2.dbo.AIRLINE_P2
    WHERE AirlineID IN (SELECT AirlineID FROM deleted);

    DELETE FROM AirDB1.dbo.AIRLINE_P1
    WHERE AirlineID IN (SELECT AirlineID FROM deleted);
END;
GO

INSERT INTO dbo.AIRLINE_View (AirlineID, IATA, ICAO, Name, Country)
VALUES (1, 'SU', 'AFL', N'Aeroflot', N'Russia');
GO

INSERT INTO dbo.AIRLINE_View (AirlineID, IATA, ICAO, Name, Country)
VALUES (2, 'AY', 'FIN', N'Finnair', N'Finland');
GO

INSERT INTO dbo.AIRLINE_View (AirlineID, IATA, ICAO, Name, Country)
VALUES (3, 'LH', 'DLH', N'Lufthansa', N'Germany');
GO

SELECT * FROM dbo.AIRLINE_View;
GO

DELETE FROM dbo.AIRLINE_View
WHERE IATA = 'LH';
GO

INSERT INTO dbo.AIRLINE_View (AirlineID, IATA, ICAO, Name, Country)
VALUES (4, 'BA', 'BAW', N'British Airways', N'United Kingdom');
GO

UPDATE dbo.AIRLINE_View
SET
    Country = N'Finland',
    Name = N'Finnair Oyj'
WHERE IATA = 'AY';
GO

SELECT * FROM AirDB1.dbo.AIRLINE_P1;
GO
SELECT * FROM AirDB2.dbo.AIRLINE_P2;
GO
