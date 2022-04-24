-- Database structure

-- Bitstreams
CREATE TABLE "bitstreams" (
	"Id"	TEXT,
	"CustomerId"	INTEGER NOT NULL,
	"SrcId"	INTEGER NOT NULL,
	"SrcOuter"	INTEGER NOT NULL,
	"SrcInner"	INTEGER NOT NULL,
	"DstId"	INTEGER NOT NULL,
	"DstOuter"	INTEGER NOT NULL,
	"DstInner"	INTEGER NOT NULL,
	"Comment"	TEXT,
	PRIMARY KEY("Id"),
    UNIQUE(SrcOuter, SrcInner),
    UNIQUE(DstOuter, DstInner)
);

-- Customers
CREATE TABLE "customers" (
	"Id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Name"	INTEGER NOT NULL,
	"OuterInterface"	INTEGER NOT NULL,
	"OuterVlan"	INTEGER NOT NULL,
	"Counter"	INTEGER NOT NULL DEFAULT 2,
    UNIQUE(OuterInterface, OuterVlan)
);

-- Flags
CREATE TABLE "flags" (
	"Key"	TEXT NOT NULL,
	"Value"	INTEGER,
	PRIMARY KEY("Key")
);

-- Default content
INSERT INTO flags (Key, Value) VALUES('IdCounter', 1000)