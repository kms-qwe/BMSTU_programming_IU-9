-- 1. AIRLINE
CREATE TABLE AIRLINE (
    AirlineID   int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    IATA        char(2)       NOT NULL,
    ICAO        char(3)       NOT NULL,
    Name        nvarchar(100) NOT NULL,
    Country     nvarchar(100) NOT NULL,

    CONSTRAINT UQ_AIRLINE_IATA UNIQUE (IATA),   -- AK 1.1
    CONSTRAINT UQ_AIRLINE_ICAO UNIQUE (ICAO)    -- AK 2.1
);
GO

-- 2. AIRPORT
CREATE TABLE AIRPORT (
    AirportID   int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    IATA        char(3)       NOT NULL,
    ICAO        char(4)       NOT NULL,
    Name        nvarchar(100) NOT NULL,
    City        nvarchar(100) NOT NULL,
    Timezone    nvarchar(50)  NOT NULL,

    CONSTRAINT UQ_AIRPORT_IATA UNIQUE (IATA),   -- AK 1.1
    CONSTRAINT UQ_AIRPORT_ICAO UNIQUE (ICAO)    -- AK 2.1
);
GO

-- 3. RUNWAY
CREATE TABLE RUNWAY (
    RunwayID    int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    AirportID   int           NOT NULL,
    Designator  char(3)       NOT NULL,
    Length      int           NOT NULL,
    Surface     nvarchar(20)  NOT NULL,

    CONSTRAINT UQ_RUNWAY_AIRPORT_DESIGNATOR
        UNIQUE (AirportID, Designator),

    CONSTRAINT FK_RUNWAY_AIRPORT
        FOREIGN KEY (AirportID)
        REFERENCES AIRPORT (AirportID)
        ON DELETE CASCADE         -- удалили аэропорт -> удалили его ВПП
        ON UPDATE CASCADE
);
GO

-- 4. AIRCRAFT
CREATE TABLE AIRCRAFT (
    AircraftID   int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    TailNumber   nvarchar(50)  NOT NULL,
    Model        nvarchar(100) NOT NULL,
    SeatCapacity int           NOT NULL,
    AirlineID    int           NOT NULL,

    CONSTRAINT UQ_AIRCRAFT_TAILNUMBER
        UNIQUE (TailNumber),                -- AK 1.1

    CONSTRAINT FK_AIRCRAFT_AIRLINE
        FOREIGN KEY (AirlineID)
        REFERENCES AIRLINE (AirlineID)
        ON DELETE CASCADE         -- удалили авиакомпанию -> удалили её самолёты
        ON UPDATE CASCADE
);
GO

-- 5. FLIGHT
CREATE TABLE FLIGHT (
    FlightID      int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    AirlineID     int           NOT NULL,   -- FK, AK 1.1
    FlightNumber  int           NOT NULL,   -- AK 1.2
    OperatingDate datetime      NOT NULL,   -- AK 1.3
    SchedDepTime  datetime      NOT NULL,
    SchedArrTime  datetime      NOT NULL,
    Status        nvarchar(20)  NOT NULL,
    DepRunwayID   int           NOT NULL,
    ArrRunwayID   int           NOT NULL,
    AircraftID    int           NOT NULL,

    CONSTRAINT UQ_FLIGHT_AIRLINE_FLIGHTNUM_OPERDATE
        UNIQUE (AirlineID, FlightNumber, OperatingDate),

    CONSTRAINT FK_FLIGHT_AIRLINE
        FOREIGN KEY (AirlineID)
        REFERENCES AIRLINE (AirlineID)
        ON DELETE NO ACTION       -- не даём удалить авиакомпанию, пока есть рейсы
        ON UPDATE CASCADE,

    CONSTRAINT FK_FLIGHT_DEPRUNWAY
        FOREIGN KEY (DepRunwayID)
        REFERENCES RUNWAY (RunwayID)
        ON DELETE NO ACTION       -- не даём удалить ВПП, пока есть рейсы
        ON UPDATE CASCADE,

    CONSTRAINT FK_FLIGHT_ARRRUNWAY
        FOREIGN KEY (ArrRunwayID)
        REFERENCES RUNWAY (RunwayID)
        ON DELETE NO ACTION
        ON UPDATE CASCADE,

    CONSTRAINT FK_FLIGHT_AIRCRAFT
        FOREIGN KEY (AircraftID)
        REFERENCES AIRCRAFT (AircraftID)
        ON DELETE NO ACTION       -- не даём удалить самолёт, пока есть рейсы
        ON UPDATE CASCADE
);
GO
