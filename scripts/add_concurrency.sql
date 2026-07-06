UPDATE courses SET slug = 'concurrency' WHERE id = '33333333-3333-3333-3333-333333333333';

-- Delete the old dummy 'concurrency' lesson
DELETE FROM lessons WHERE slug = 'concurrency' AND course_id = '33333333-3333-3333-3333-333333333333';

-- Insert the 25 real lessons
INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '01-Introduction', '01: Introduction', 'See markdown file', 1),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '02-Why-Concurrency', '02: Why Concurrency', 'See markdown file', 2),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '03-Concurrency-vs-Parallelism', '03: Concurrency vs Parallelism', 'See markdown file', 3),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '04-Process-vs-Thread-vs-Goroutine', '04: Process vs Thread vs Goroutine', 'See markdown file', 4),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '05-Go-Runtime', '05: Go Runtime', 'See markdown file', 5),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '06-Go-Scheduler', '06: Go Scheduler', 'See markdown file', 6),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '07-GPM-Model', '07: GPM Model', 'See markdown file', 7),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '08-Goroutines', '08: Goroutines', 'See markdown file', 8),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '09-WaitGroup', '09: WaitGroup', 'See markdown file', 9),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '10-Channels', '10: Channels', 'See markdown file', 10),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '11-Buffered-Channels', '11: Buffered Channels', 'See markdown file', 11),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '12-Unbuffered-Channels', '12: Unbuffered Channels', 'See markdown file', 12),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '13-Channel-Directions', '13: Channel Directions', 'See markdown file', 13),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '14-Channel-Closing', '14: Channel Closing', 'See markdown file', 14),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '15-Range', '15: Range', 'See markdown file', 15),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '16-Select', '16: Select', 'See markdown file', 16),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '17-Timers', '17: Timers', 'See markdown file', 17),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '18-Tickers', '18: Tickers', 'See markdown file', 18),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '19-Context', '19: Context', 'See markdown file', 19),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '20-Cancellation', '20: Cancellation', 'See markdown file', 20),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '21-Mutex', '21: Mutex', 'See markdown file', 21),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '22-RWMutex', '22: RWMutex', 'See markdown file', 22),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '23-Atomic', '23: Atomic', 'See markdown file', 23),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '24-sync.Once', '24: sync.Once', 'See markdown file', 24),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '25-sync.Map', '25: sync.Map', 'See markdown file', 25)
ON CONFLICT DO NOTHING;
