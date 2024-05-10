-- file: 20-create-database.sql
CREATE DATABASE commands;
GRANT ALL PRIVILEGES ON DATABASE commands TO program;

\c commands;

CREATE TABLE command (
    id                SERIAL PRIMARY KEY,
    command_uuid uuid UNIQUE NOT NULL,
    description       VARCHAR(255) NOT NULL,
    script            TEXT NOT NULL,
    status            VARCHAR(20) DEFAULT 'UNKNOWN'
        CHECK (status IN ('EXECUTING', 'EXECUTED', 'STOPPED', 'FAILED', 'UNKNOWN')),
    output            TEXT NOT NULL
);

GRANT ALL ON ALL TABLES IN SCHEMA public TO program;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO program;


-- INSERT INTO command (command_uuid, description, script, status, output) VALUES ('2285996a-4f8d-40f6-8f2a-07c35daf4d94', 'Команда через 3 секунды выводит Hello World', 'sleep 3; echo Hello World', 'EXECUTED', 'Hello World\n');













-- INSERT INTO banner (feature_id, content, is_active, created_at, updated_at) VALUES (3, '{ "title": "Alice", "text": "alice", "url": "alice.com"}', true, '2024-04-10T10:03:05', '2024-04-14T12:03:05');
-- INSERT INTO banner_tag VALUES (2, 1);

-- INSERT INTO banner (feature_id, content, is_active, created_at, updated_at) VALUES (3, '{ "title": "AliceGood", "text": "alicegood", "url": "alicegood.com"}', false, '2024-04-10T10:03:05', '2024-04-14T12:03:05');
-- INSERT INTO banner_tag VALUES (3, 1);

-- SELECT b.id, b.content, b.feature_id, array_agg(bt.tag_id) as tag_ids
-- FROM banner b
-- JOIN banner_tag bt ON b.id = bt.banner_id
-- WHERE b.feature_id = 8
-- GROUP BY b.id, b.feature_id
-- HAVING bool_or(bt.tag_id = 10)
-- LIMIT 1
-- OFFSET 1;

-- SELECT b.id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at
-- FROM banner b
-- JOIN banner_tag bt ON b.id = bt.banner_id
-- WHERE b.feature_id = 8 AND bt.tag_id = 10;

-- SELECT b.content
-- FROM banner b
-- JOIN banner_tag bt ON b.id = bt.banner_id
-- WHERE b.feature_id = 8 AND bt.tag_id = 10;

-- UPDATE banner SET feature_id = 3, content = '{ "title": "Alice", "text": "alice", "url": "alice.com"}', is_active = true, updated_at = '2024-04-12T23:00:00' WHERE id = 2;

