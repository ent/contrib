# todo example with PULIDs (Prefixed ULIDs).

`PULID`s are an identifier encoding scheme that builds upon the excellent [ULID](Universally
Unique Lexicographically Sortable Identifier) scheme. Prefixes should maintain compatibility with
base32. If constrained to 2 characters for encoding the entity type, the number of types of entities that 
can be referenced with this scheme is `2^32` (1024).


