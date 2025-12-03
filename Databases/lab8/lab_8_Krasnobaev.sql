USE master;
GO

IF DB_ID(N'Lab8') IS NOT NULL
BEGIN
    ALTER DATABASE Lab8 SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
    DROP DATABASE Lab8;
END;
GO

CREATE DATABASE Lab8
ON ( NAME = Lab8_dat,
    FILENAME = '/var/opt/mssql/data/Lab8_dat.mdf',
    SIZE = 10MB,
    MAXSIZE = UNLIMITED,
    FILEGROWTH = 5%
)
LOG ON (NAME = Lab8_log,
    FILENAME = '/var/opt/mssql/data/Lab8_log.ldf',
    SIZE = 5MB,
    MAXSIZE = 25MB,
    FILEGROWTH = 5MB
);
GO

USE Lab8;
GO

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

INSERT INTO dbo.AIRLINE (IATA, ICAO, Name, Country)
VALUES
    ('SU', 'AFL', N'Aeroflot',        N'Russia'),
    ('S7', 'SBI', N'S7 Airlines',     N'Russia'),
    ('LH', 'DLH', N'Lufthansa',       N'Germany'),
    ('AF', 'AFR', N'Air France',      N'France'),
    ('BA', 'BAW', N'British Airways', N'United Kingdom');
GO

DROP PROCEDURE IF EXISTS dbo.GetAirlinesByCountryLength;
GO

CREATE PROCEDURE dbo.GetAirlinesByCountryLength
    @cursor        CURSOR VARYING OUTPUT,
    @CountryLength INT
AS
BEGIN
    SET @cursor = CURSOR FORWARD_ONLY STATIC FOR
    SELECT AirlineID, Name, Country, IATA, ICAO
    FROM dbo.AIRLINE
    WHERE LEN(Country) = @CountryLength;

    OPEN @cursor;
END;
GO

DROP FUNCTION IF EXISTS dbo.TransformAirlineName;
GO

CREATE FUNCTION dbo.TransformAirlineName(@Name NVARCHAR(100))
RETURNS NVARCHAR(100)
AS
BEGIN
    RETURN REVERSE(@Name);
END;
GO

DROP PROCEDURE IF EXISTS dbo.GetTransformedAirlines;
GO

CREATE PROCEDURE dbo.GetTransformedAirlines
    @cursor        CURSOR VARYING OUTPUT,
    @CountryLength INT
AS
BEGIN
    SET @cursor = CURSOR FORWARD_ONLY STATIC FOR
    SELECT
        AirlineID,
        dbo.TransformAirlineName(Name) AS TransformedName,
        Country,
        IATA,
        ICAO
    FROM dbo.AIRLINE
    WHERE LEN(Country) = @CountryLength;

    OPEN @cursor;
END;
GO

DROP FUNCTION IF EXISTS dbo.IsRussianAirline;
GO

CREATE FUNCTION dbo.IsRussianAirline(@Country NVARCHAR(100))
RETURNS BIT
AS
BEGIN
    RETURN CASE WHEN @Country = N'Russia' THEN 1 ELSE 0 END;
END;
GO

DROP PROCEDURE IF EXISTS dbo.GetRussianAirlinesData;
GO

CREATE PROCEDURE dbo.GetRussianAirlinesData
AS
BEGIN
    DECLARE @cursor CURSOR;

    EXEC dbo.GetAirlinesByCountryLength
        @cursor        = @cursor OUTPUT,
        @CountryLength = 6;

    DECLARE
        @AirlineID INT,
        @Name      NVARCHAR(100),
        @Country   NVARCHAR(100),
        @IATA      CHAR(2),
        @ICAO      CHAR(3);

    FETCH NEXT FROM @cursor
        INTO @AirlineID, @Name, @Country, @IATA, @ICAO;

    PRINT 'Results for dbo.GetRussianAirlinesData:';

    WHILE @@FETCH_STATUS = 0
    BEGIN
        IF dbo.IsRussianAirline(@Country) = 1
        BEGIN
            SELECT
                @AirlineID AS AirlineID,
                @Name      AS AirlineName,
                @Country   AS Country,
                @IATA      AS IATA,
                @ICAO      AS ICAO;
        END;

        FETCH NEXT FROM @cursor
            INTO @AirlineID, @Name, @Country, @IATA, @ICAO;
    END;

    CLOSE @cursor;
    DEALLOCATE @cursor;
END;
GO

DROP FUNCTION IF EXISTS dbo.FilteredTransformedAirlines;
GO

CREATE FUNCTION dbo.FilteredTransformedAirlines()
RETURNS TABLE
AS
RETURN
(
    SELECT
        AirlineID,
        dbo.TransformAirlineName(Name) AS Name,
        Country,
        IATA,
        ICAO
    FROM dbo.AIRLINE
    WHERE Country = N'Russia'
);
GO

-- CREATE FUNCTION dbo.FilteredTransformedAirlines()
-- RETURNS @ResultTable TABLE
-- (
--     AirlineID INT,
--     Name NVARCHAR(100),
--     Country NVARCHAR(100),
--     IATA CHAR(2),
--     ICAO CHAR(3)
-- )
-- AS
-- BEGIN
--     INSERT INTO @ResultTable (AirlineID, Name, Country, IATA, ICAO)
--     SELECT
--         AirlineID,
--         dbo.TransformAirlineName(Name),
--         Country,
--         IATA,
--         ICAO
--     FROM dbo.AIRLINE
--     WHERE Country = N'Russia';

--     RETURN;
-- END;
-- GO

DROP PROCEDURE IF EXISTS dbo.UseFilteredTransformedAirlines;
GO

CREATE PROCEDURE dbo.UseFilteredTransformedAirlines
AS
BEGIN
    SELECT *
    FROM dbo.FilteredTransformedAirlines();
END;
GO

DECLARE @cursor1 CURSOR;
DECLARE
    @AirlineID1 INT,
    @Name1      NVARCHAR(100),
    @Country1   NVARCHAR(100),
    @IATA1      CHAR(2),
    @ICAO1      CHAR(3);

EXEC dbo.GetAirlinesByCountryLength
    @cursor        = @cursor1 OUTPUT,
    @CountryLength = 6;

FETCH NEXT FROM @cursor1 INTO @AirlineID1, @Name1, @Country1, @IATA1, @ICAO1;
WHILE @@FETCH_STATUS = 0
BEGIN
    SELECT
        @AirlineID1 AS AirlineID,
        @Name1      AS AirlineName,
        @Country1   AS Country,
        @IATA1      AS IATA,
        @ICAO1      AS ICAO;

    FETCH NEXT FROM @cursor1 INTO @AirlineID1, @Name1, @Country1, @IATA1, @ICAO1;
END;
CLOSE @cursor1;
DEALLOCATE @cursor1;
GO

DECLARE @cursor2 CURSOR;
DECLARE
    @AirlineID2 INT,
    @Name2      NVARCHAR(100),
    @Country2   NVARCHAR(100),
    @IATA2      CHAR(2),
    @ICAO2      CHAR(3);

EXEC dbo.GetTransformedAirlines
    @cursor        = @cursor2 OUTPUT,
    @CountryLength = 6;

FETCH NEXT FROM @cursor2 INTO @AirlineID2, @Name2, @Country2, @IATA2, @ICAO2;
WHILE @@FETCH_STATUS = 0
BEGIN
    SELECT
        @AirlineID2 AS AirlineID,
        @Name2      AS TransformedName,
        @Country2   AS Country,
        @IATA2      AS IATA,
        @ICAO2      AS ICAO;

    FETCH NEXT FROM @cursor2 INTO @AirlineID2, @Name2, @Country2, @IATA2, @ICAO2;
END;
CLOSE @cursor2;
DEALLOCATE @cursor2;
GO

EXEC dbo.GetRussianAirlinesData;
GO

EXEC dbo.UseFilteredTransformedAirlines;
GO
