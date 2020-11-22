/*
   When inserting documents, the column definition is omitted. The user can provide JSON literals as documents.
 */
INSERT INTO movies
VALUES ('{"title":"hello","year":2938}'),
       ('{"title":"foo","year":2938}');

/*
   When using placeholders, (?) automatically means list of documents. The user may provide a single item or a slice
   of items.
 */
INSERT INTO movies
VALUES (?);
