/* psql -d DATABASE < seed.sql */

INSERT INTO users(id, username) VALUES(60420644,'Yuzhay');
INSERT INTO currencies(id, code, rate) VALUES(1,'USD', 1.1), (2, 'EUR', 1.2), (3, 'GBP', 1.3);
INSERT INTO subscriptions(id, user_id, currency_id) VALUES(1,60420644,1), (2,60420644,2), (3,60420644,3);