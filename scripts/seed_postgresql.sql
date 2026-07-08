
-- Seed PostgreSQL Course
INSERT INTO courses (id, slug, title, description, difficulty) VALUES
('c0000000-0000-0000-0000-000000000002', 'postgresql', 'PostgreSQL Mastery', 'From Basic to Advanced to Expert. Learn PostgreSQL internals, real-world backend patterns, and Go integration.', 'expert')
ON CONFLICT (id) DO UPDATE SET slug = EXCLUDED.slug, title = EXCLUDED.title, description = EXCLUDED.description, difficulty = EXCLUDED.difficulty;

DELETE FROM lessons WHERE course_id = 'c0000000-0000-0000-0000-000000000002';

-- Seed PostgreSQL Lessons
INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
('30000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', '001-level-1-postgresql-fundamentals', 'Level 1 - PostgreSQL Fundamentals', 'See markdown file', 1),
('30000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000002', '002-level-2-crud-operations', 'Level 2 - CRUD Operations', 'See markdown file', 2),
('30000000-0000-0000-0000-000000000003', 'c0000000-0000-0000-0000-000000000002', '003-level-3-constraints', 'Level 3 - Constraints', 'See markdown file', 3),
('30000000-0000-0000-0000-000000000004', 'c0000000-0000-0000-0000-000000000002', '004-level-4-querying', 'Level 4 - Querying', 'See markdown file', 4),
('30000000-0000-0000-0000-000000000005', 'c0000000-0000-0000-0000-000000000002', '005-level-5-joins', 'Level 5 - Joins', 'See markdown file', 5),
('30000000-0000-0000-0000-000000000006', 'c0000000-0000-0000-0000-000000000002', '006-level-6-aggregate-functions', 'Level 6 - Aggregate Functions', 'See markdown file', 6),
('30000000-0000-0000-0000-000000000007', 'c0000000-0000-0000-0000-000000000002', '007-level-7-built-in-functions', 'Level 7 - Built-in Functions', 'See markdown file', 7),
('30000000-0000-0000-0000-000000000008', 'c0000000-0000-0000-0000-000000000002', '008-level-8-views', 'Level 8 - Views', 'See markdown file', 8),
('30000000-0000-0000-0000-000000000009', 'c0000000-0000-0000-0000-000000000002', '009-level-9-indexing', 'Level 9 - Indexing', 'See markdown file', 9),
('30000000-0000-0000-0000-000000000010', 'c0000000-0000-0000-0000-000000000002', '010-level-10-transactions', 'Level 10 - Transactions', 'See markdown file', 10),
('30000000-0000-0000-0000-000000000011', 'c0000000-0000-0000-0000-000000000002', '011-level-11-normalization', 'Level 11 - Normalization', 'See markdown file', 11),
('30000000-0000-0000-0000-000000000012', 'c0000000-0000-0000-0000-000000000002', '012-level-12-advanced-sql', 'Level 12 - Advanced SQL', 'See markdown file', 12),
('30000000-0000-0000-0000-000000000013', 'c0000000-0000-0000-0000-000000000002', '013-level-13-json-jsonb', 'Level 13 - JSON & JSONB', 'See markdown file', 13),
('30000000-0000-0000-0000-000000000014', 'c0000000-0000-0000-0000-000000000002', '014-level-14-arrays', 'Level 14 - Arrays', 'See markdown file', 14),
('30000000-0000-0000-0000-000000000015', 'c0000000-0000-0000-0000-000000000002', '015-level-15-uuid', 'Level 15 - UUID', 'See markdown file', 15),
('30000000-0000-0000-0000-000000000016', 'c0000000-0000-0000-0000-000000000002', '016-level-16-sequences', 'Level 16 - Sequences', 'See markdown file', 16),
('30000000-0000-0000-0000-000000000017', 'c0000000-0000-0000-0000-000000000002', '017-level-17-stored-procedures', 'Level 17 - Stored Procedures', 'See markdown file', 17),
('30000000-0000-0000-0000-000000000018', 'c0000000-0000-0000-0000-000000000002', '018-level-18-triggers', 'Level 18 - Triggers', 'See markdown file', 18),
('30000000-0000-0000-0000-000000000019', 'c0000000-0000-0000-0000-000000000002', '019-level-19-plpgsql', 'Level 19 - PL/pgSQL', 'See markdown file', 19),
('30000000-0000-0000-0000-000000000020', 'c0000000-0000-0000-0000-000000000002', '020-level-20-security', 'Level 20 - Security', 'See markdown file', 20),
('30000000-0000-0000-0000-000000000021', 'c0000000-0000-0000-0000-000000000002', '021-level-21-backup-restore', 'Level 21 - Backup & Restore', 'See markdown file', 21),
('30000000-0000-0000-0000-000000000022', 'c0000000-0000-0000-0000-000000000002', '022-level-22-performance-tuning', 'Level 22 - Performance Tuning', 'See markdown file', 22),
('30000000-0000-0000-0000-000000000023', 'c0000000-0000-0000-0000-000000000002', '023-level-23-partitioning', 'Level 23 - Partitioning', 'See markdown file', 23),
('30000000-0000-0000-0000-000000000024', 'c0000000-0000-0000-0000-000000000002', '024-level-24-replication', 'Level 24 - Replication', 'See markdown file', 24),
('30000000-0000-0000-0000-000000000025', 'c0000000-0000-0000-0000-000000000002', '025-level-25-high-availability', 'Level 25 - High Availability', 'See markdown file', 25),
('30000000-0000-0000-0000-000000000026', 'c0000000-0000-0000-0000-000000000002', '026-level-26-full-text-search', 'Level 26 - Full-Text Search', 'See markdown file', 26),
('30000000-0000-0000-0000-000000000027', 'c0000000-0000-0000-0000-000000000002', '027-level-27-extensions', 'Level 27 - Extensions', 'See markdown file', 27),
('30000000-0000-0000-0000-000000000028', 'c0000000-0000-0000-0000-000000000002', '028-level-28-monitoring', 'Level 28 - Monitoring', 'See markdown file', 28),
('30000000-0000-0000-0000-000000000029', 'c0000000-0000-0000-0000-000000000002', '029-level-29-concurrency', 'Level 29 - Concurrency', 'See markdown file', 29),
('30000000-0000-0000-0000-000000000030', 'c0000000-0000-0000-0000-000000000002', '030-level-30-advanced-storage', 'Level 30 - Advanced Storage', 'See markdown file', 30),
('30000000-0000-0000-0000-000000000031', 'c0000000-0000-0000-0000-000000000002', '031-level-31-postgresql-internals', 'Level 31 - PostgreSQL Internals', 'See markdown file', 31),
('30000000-0000-0000-0000-000000000032', 'c0000000-0000-0000-0000-000000000002', '032-level-32-distributed-postgresql', 'Level 32 - Distributed PostgreSQL', 'See markdown file', 32),
('30000000-0000-0000-0000-000000000033', 'c0000000-0000-0000-0000-000000000002', '033-level-33-postgresql-with-go', 'Level 33 - PostgreSQL with Go', 'See markdown file', 33),
('30000000-0000-0000-0000-000000000034', 'c0000000-0000-0000-0000-000000000002', '034-level-34-interview-topics', 'Level 34 - Interview Topics', 'See markdown file', 34),
('30000000-0000-0000-0000-000000000035', 'c0000000-0000-0000-0000-000000000002', '035-level-35-real-world-backend-patterns', 'Level 35 - Real-World Backend Patterns', 'See markdown file', 35)
;
