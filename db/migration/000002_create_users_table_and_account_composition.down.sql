DROP INDEX owner_currency;

ALTER TABLE accounts DROP CONSTRAINT accounts_owner_fkey;

DROP TABLE IF EXISTS users;
