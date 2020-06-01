CREATE TYPE activity_type AS ENUM (
  'genesis',
  'transfer',
  'deposit',
  'unbond',
  'withdraw',
  'nodejoin',
  'unjail',
  'reward',
  'jail',
  'slash',
  'nodekicked'
);