-- Add Networking Course
INSERT INTO courses (id, slug, title, description, difficulty) 
VALUES ('11111111-1111-1111-1111-111111111111', 'networking', 'Networking & Socket Programming', 'Learn TCP, UDP, WebSockets, and Reverse Proxies in Go.', 'advanced')
ON CONFLICT (id) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description;

-- Add Profiling Course
INSERT INTO courses (id, slug, title, description, difficulty) 
VALUES ('22222222-2222-2222-2222-222222222222', 'profiling', 'Profiling in Go', 'Master pprof, execution tracing, and memory leak detection.', 'expert')
ON CONFLICT (id) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description;

-- Update Performance Engineering to include the new lessons (assuming id is aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa)
-- Wait, performance-engineering already exists with ID 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'.

-- Insert Networking Lessons
INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
('33333333-3333-3333-3333-000000000001', '11111111-1111-1111-1111-111111111111', '001-tcp-sockets-in-go', 'TCP Sockets in Go', 'See markdown file', 1),
('33333333-3333-3333-3333-000000000002', '11111111-1111-1111-1111-111111111111', '002-udp-websockets', 'UDP and WebSockets', 'See markdown file', 2),
('33333333-3333-3333-3333-000000000003', '11111111-1111-1111-1111-111111111111', '003-reverse-proxy', 'Building a Reverse Proxy', 'See markdown file', 3)
ON CONFLICT (id) DO UPDATE SET slug = EXCLUDED.slug, title = EXCLUDED.title, content = EXCLUDED.content;

-- Insert Profiling Lessons
INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
('44444444-4444-4444-4444-000000000001', '22222222-2222-2222-2222-222222222222', '001-pprof-basics', 'Introduction to pprof', 'See markdown file', 1),
('44444444-4444-4444-4444-000000000002', '22222222-2222-2222-2222-222222222222', '002-flame-graphs-and-tracing', 'Flame Graphs and Tracing', 'See markdown file', 2),
('44444444-4444-4444-4444-000000000003', '22222222-2222-2222-2222-222222222222', '003-memory-leaks-and-pprof', 'Memory Leaks and pprof', 'See markdown file', 3)
ON CONFLICT (id) DO UPDATE SET slug = EXCLUDED.slug, title = EXCLUDED.title, content = EXCLUDED.content;

-- Insert Performance Engineering Lessons (adding the new ones, offsetting order index to be at the end, or replacing existing)
-- Since the existing course 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa' has lessons up to index 6. Let's add them at 7, 8.
-- Wait, the new ones I created were in a folder named 'performance', NOT 'performance-engineering'.
-- I will create a new course called 'performance' to perfectly match the folder name so the markdown files load correctly!
INSERT INTO courses (id, slug, title, description, difficulty) 
VALUES ('55555555-5555-5555-5555-555555555555', 'performance', 'High Performance Go', 'Mechanical Sympathy, Sync Pools, and Lock-Free coding.', 'expert')
ON CONFLICT (id) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description;

INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
('66666666-6666-6666-6666-000000000001', '55555555-5555-5555-5555-555555555555', '001-memory-allocations', 'Mechanical Sympathy & Memory', 'See markdown file', 1),
('66666666-6666-6666-6666-000000000002', '55555555-5555-5555-5555-555555555555', '002-sync-pool', 'Object Reuse with sync.Pool', 'See markdown file', 2),
('66666666-6666-6666-6666-000000000003', '55555555-5555-5555-5555-555555555555', '003-concurrency-optimizations', 'Concurrency Optimizations', 'See markdown file', 3)
ON CONFLICT (id) DO UPDATE SET slug = EXCLUDED.slug, title = EXCLUDED.title, content = EXCLUDED.content;
