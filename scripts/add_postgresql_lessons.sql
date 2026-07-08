INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
('30000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', '001-architecture-and-storage', 'Architecture and Storage', 'See markdown file', 1),
('30000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000002', '002-transaction-isolation-levels', 'Transaction Isolation Levels', 'See markdown file', 2),
('30000000-0000-0000-0000-000000000003', 'c0000000-0000-0000-0000-000000000002', '003-advanced-indexing', 'Advanced Indexing', 'See markdown file', 3),
('30000000-0000-0000-0000-000000000004', 'c0000000-0000-0000-0000-000000000002', '004-jsonb-and-document-storage', 'JSONB and Document Storage', 'See markdown file', 4),
('30000000-0000-0000-0000-000000000005', 'c0000000-0000-0000-0000-000000000002', '005-connection-pooling-pgbouncer', 'Connection Pooling PgBouncer', 'See markdown file', 5),
('30000000-0000-0000-0000-000000000006', 'c0000000-0000-0000-0000-000000000002', '006-partitioning-and-sharding', 'Partitioning and Sharding', 'See markdown file', 6)
ON CONFLICT (id) DO UPDATE SET slug = EXCLUDED.slug, title = EXCLUDED.title, content = EXCLUDED.content, order_index = EXCLUDED.order_index;
