USE master;
GO

IF DB_ID(N'Lab9') IS NOT NULL
BEGIN
    ALTER DATABASE Lab9 SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
    DROP DATABASE Lab9;
END;
GO

CREATE DATABASE Lab9
ON (
    NAME = Lab9_dat,
    FILENAME = '/var/opt/mssql/data/Lab9_dat.mdf',
    SIZE = 10MB,
    MAXSIZE = UNLIMITED,
    FILEGROWTH = 5%
)
LOG ON (
    NAME = Lab9_log,
    FILENAME = '/var/opt/mssql/data/Lab9_log.ldf',
    SIZE = 5MB,
    MAXSIZE = 5GB,
    FILEGROWTH = 5MB
);
GO

USE Lab9;
GO

DROP TABLE IF EXISTS dbo.AIRLINE_CEO;
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

CREATE TABLE dbo.AIRLINE_CEO (
    AirlineID      INT           NOT NULL PRIMARY KEY,
    CEOFirstName   NVARCHAR(50)  NOT NULL,
    CEOLastName    NVARCHAR(50)  NOT NULL,
    YearsInService INT           NOT NULL CHECK (YearsInService >= 0),
    CONSTRAINT FK_AIRLINE_CEO_AIRLINE FOREIGN KEY (AirlineID)
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

INSERT INTO dbo.AIRLINE_CEO (AirlineID, CEOFirstName, CEOLastName, YearsInService)
VALUES
    (1, N'Sergey',  N'Alexandrov', 5),
    (2, N'Igor',    N'Petrov',     3),
    (3, N'Heinrich',N'Meyer',      4),
    (4, N'Claude',  N'Dupont',     6),
    (5, N'John',    N'Smith',      2);
GO

DROP TRIGGER IF EXISTS dbo.trg_Insert_AIRLINE;
GO

CREATE TRIGGER trg_Insert_AIRLINE
ON dbo.AIRLINE
AFTER INSERT
AS
BEGIN
    SELECT * FROM inserted;
END;
GO

DROP TRIGGER IF EXISTS dbo.trg_Update_AIRLINE;
GO

CREATE TRIGGER trg_Update_AIRLINE
ON dbo.AIRLINE
AFTER UPDATE
AS
BEGIN
    IF EXISTS (SELECT 1 FROM inserted WHERE LEN(Name) < 3)
    BEGIN
        RAISERROR('Airline name must be at least 3 characters.', 16, 1);
        ROLLBACK TRANSACTION;
        RETURN;
    END;

    SELECT * FROM inserted;
END;
GO

DROP TRIGGER IF EXISTS dbo.trg_Delete_AIRLINE;
GO

CREATE TRIGGER trg_Delete_AIRLINE
ON dbo.AIRLINE
AFTER DELETE
AS
BEGIN
    SELECT * FROM deleted;
END;
GO

DROP VIEW IF EXISTS dbo.vw_AirlineWithCeo;
GO

CREATE VIEW dbo.vw_AirlineWithCeo
AS
    SELECT
        al.AirlineID,
        al.IATA,
        al.ICAO,
        al.Name           AS AirlineName,
        al.Country,
        c.CEOFirstName,
        c.CEOLastName,
        c.YearsInService
    FROM dbo.AIRLINE al
    INNER JOIN dbo.AIRLINE_CEO c ON al.AirlineID = c.AirlineID;
GO

DROP TRIGGER IF EXISTS dbo.trg_Insert_vw_AirlineWithCeo;
GO

CREATE TRIGGER trg_Insert_vw_AirlineWithCeo
ON dbo.vw_AirlineWithCeo
INSTEAD OF INSERT
AS
BEGIN

    IF EXISTS (SELECT 1 FROM inserted WHERE AirlineID IS NOT NULL)
    BEGIN
        RAISERROR('AirlineID cannot be specified on insert.', 16, 1);
        ROLLBACK TRANSACTION;
        RETURN;
    END;

    IF EXISTS (
        SELECT 1
        FROM inserted i
        JOIN dbo.AIRLINE al
            ON al.IATA = i.IATA
            OR al.ICAO = i.ICAO
    )
    BEGIN
        RAISERROR('Airline with this IATA or ICAO already exists.', 16, 1);
        ROLLBACK TRANSACTION;
        RETURN;
    END;

    INSERT INTO dbo.AIRLINE (IATA, ICAO, Name, Country)
    SELECT
        IATA,
        ICAO,
        AirlineName,
        Country
    FROM inserted;

    INSERT INTO dbo.AIRLINE_CEO (AirlineID, CEOFirstName, CEOLastName, YearsInService)
    SELECT
        al.AirlineID,
        i.CEOFirstName,
        i.CEOLastName,
        i.YearsInService
    FROM inserted i
    INNER JOIN dbo.AIRLINE al
        ON al.IATA = i.IATA
       AND al.ICAO = i.ICAO;
END;
GO

DROP TRIGGER IF EXISTS dbo.trg_Update_vw_AirlineWithCeo;
GO

CREATE TRIGGER trg_Update_vw_AirlineWithCeo
ON dbo.vw_AirlineWithCeo
INSTEAD OF UPDATE
AS
BEGIN

    IF UPDATE(AirlineID)
    BEGIN
        RAISERROR('Нельзя изменять столбец AirlineID через это представление.', 16, 1);
        ROLLBACK TRANSACTION;
        RETURN;
    END;

    IF EXISTS (SELECT 1 FROM inserted WHERE LEN(CEOLastName) < 3)
    BEGIN
        RAISERROR('CEO last name must be at least 3 characters.', 16, 1);
        ROLLBACK TRANSACTION;
        RETURN;
    END;

    WITH UpdatedAirline AS
    (
        SELECT
            al.AirlineID,
            i.IATA,
            i.ICAO,
            i.AirlineName,
            i.Country
        FROM dbo.AIRLINE al
        INNER JOIN inserted i ON al.AirlineID = i.AirlineID
    )
    UPDATE al
    SET
        IATA    = ua.IATA,
        ICAO    = ua.ICAO,
        Name    = ua.AirlineName,
        Country = ua.Country
    FROM dbo.AIRLINE al
    INNER JOIN UpdatedAirline ua ON al.AirlineID = ua.AirlineID;

    WITH UpdatedCeo AS
    (
        SELECT
            c.AirlineID,
            i.CEOFirstName,
            i.CEOLastName,
            i.YearsInService
        FROM dbo.AIRLINE_CEO c
        INNER JOIN inserted i ON c.AirlineID = i.AirlineID
    )
    UPDATE c
    SET
        CEOFirstName   = uc.CEOFirstName,
        CEOLastName    = uc.CEOLastName,
        YearsInService = uc.YearsInService
    FROM dbo.AIRLINE_CEO c
    INNER JOIN UpdatedCeo uc ON c.AirlineID = uc.AirlineID;
END;
GO

DROP TRIGGER IF EXISTS dbo.trg_Delete_vw_AirlineWithCeo;
GO

CREATE TRIGGER trg_Delete_vw_AirlineWithCeo
ON dbo.vw_AirlineWithCeo
INSTEAD OF DELETE
AS
BEGIN
    DELETE FROM dbo.AIRLINE_CEO
    WHERE AirlineID IN (SELECT AirlineID FROM deleted);

    DELETE FROM dbo.AIRLINE
    WHERE AirlineID IN (SELECT AirlineID FROM deleted);
END;
GO

INSERT INTO dbo.vw_AirlineWithCeo
    (IATA, ICAO, AirlineName, Country, CEOFirstName, CEOLastName, YearsInService)
VALUES
    ('UT', 'UTA', N'Utair', N'Russia', N'Andrey', N'Martynov', 5);
GO

UPDATE dbo.vw_AirlineWithCeo
SET
    CEOFirstName   = N'Andrey',
    CEOLastName    = N'Ivanov',
    YearsInService = 6
WHERE IATA = 'UT';
GO

DELETE FROM dbo.vw_AirlineWithCeo
WHERE IATA = 'UT';
GO

BEGIN TRY
    UPDATE dbo.AIRLINE
    SET Name = N'AA'
    WHERE AirlineID = 1;
END TRY
BEGIN CATCH
END CATCH;
GO

BEGIN TRY
    UPDATE dbo.vw_AirlineWithCeo
    SET CEOLastName = N'A'
    WHERE AirlineID = 2;
END TRY
BEGIN CATCH
END CATCH;
GO

SELECT * FROM dbo.AIRLINE;
SELECT * FROM dbo.AIRLINE_CEO;
SELECT * FROM dbo.vw_AirlineWithCeo;
update dbo.vw_AirlineWithCeo set IATA=IATA;
update dbo.vw_AirlineWithCeo set AirlineID = AirlineID+ 1;
update dbo.AIRLINE set Name = UPPER(Name);
SELECT * FROM dbo.vw_AirlineWithCeo;
GO
