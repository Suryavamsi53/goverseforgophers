-- Insert the two new courses
INSERT INTO courses (id, slug, title, description, difficulty) VALUES 
('77777777-7777-7777-7777-777777777777', 'design-patterns', 'Design Patterns in Go', 'Master classic software design patterns and idiomatic Go design structures.', 'intermediate'),
('88888888-8888-8888-8888-888888888888', 'distributed-systems', 'Distributed Systems in Go', 'Build reliable, resilient, consistent, and observable cloud-native systems in Go.', 'advanced')
ON CONFLICT (slug) DO NOTHING;

-- Insert remaining Concurrency lessons (26 to 41)
INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '26-sync.Cond', '26: sync.Cond', 'See markdown file', 26),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '27-errgroup', '27: errgroup', 'See markdown file', 27),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '28-Race-Conditions', '28: Race Conditions', 'See markdown file', 28),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '29-Deadlocks', '29: Deadlocks', 'See markdown file', 29),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '30-Starvation', '30: Starvation', 'See markdown file', 30),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '31-Livelock', '31: Livelock', 'See markdown file', 31),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '32-Worker-Pool', '32: Worker Pool', 'See markdown file', 32),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '33-Fan-In', '33: Fan-In', 'See markdown file', 33),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '34-Fan-Out', '34: Fan-Out', 'See markdown file', 34),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '35-Pipeline', '35: Pipeline', 'See markdown file', 35),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '36-Semaphore', '36: Semaphore', 'See markdown file', 36),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '37-Rate-Limiter', '37: Rate Limiter', 'See markdown file', 37),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '38-Performance', '38: Performance', 'See markdown file', 38),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '39-Error-Handling', '39: Error Handling', 'See markdown file', 39),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '40-Testing-Concurrency', '40: Testing Concurrency', 'See markdown file', 40),
(gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '41-Conclusion', '41: Conclusion', 'See markdown file', 41)
ON CONFLICT (slug) DO NOTHING;

-- Insert Design Patterns lessons (1 to 17)
INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '01-Functional-Options', '01: Functional Options', 'See markdown file', 1),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '02-Accept-Interfaces-Return-Structs', '02: Accept Interfaces Return Structs', 'See markdown file', 2),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '03-Context-Propagation', '03: Context Propagation', 'See markdown file', 3),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '04-Error-Wrapping-and-Typing', '04: Error Wrapping and Typing', 'See markdown file', 4),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '05-Factory-Method', '05: Factory Method', 'See markdown file', 5),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '06-Builder', '06: Builder', 'See markdown file', 6),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '07-Singleton', '07: Singleton', 'See markdown file', 7),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '08-Object-Pool', '08: Object Pool', 'See markdown file', 8),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '09-Adapter', '09: Adapter', 'See markdown file', 9),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '10-Decorator', '10: Decorator', 'See markdown file', 10),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '11-Facade', '11: Facade', 'See markdown file', 11),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '12-Proxy', '12: Proxy', 'See markdown file', 12),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '13-Strategy', '13: Strategy', 'See markdown file', 13),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '14-Observer', '14: Observer', 'See markdown file', 14),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '15-Command', '15: Command', 'See markdown file', 15),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '16-State', '16: State', 'See markdown file', 16),
(gen_random_uuid(), '77777777-7777-7777-7777-777777777777', '17-Chain-of-Responsibility', '17: Chain of Responsibility', 'See markdown file', 17)
ON CONFLICT (slug) DO NOTHING;

-- Insert Distributed Systems lessons (1 to 15)
INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '01-CAP-Theorem', '01: CAP Theorem', 'See markdown file', 1),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '02-Fallacies-of-Distributed-Computing', '02: Fallacies of Distributed Computing', 'See markdown file', 2),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '03-Time-and-Clocks', '03: Time and Clocks', 'See markdown file', 3),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '04-RPC-vs-REST', '04: RPC vs REST', 'See markdown file', 4),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '05-gRPC-and-Protobuf', '05: gRPC and Protobuf', 'See markdown file', 5),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '06-Message-Queues', '06: Message Queues', 'See markdown file', 6),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '07-Circuit-Breaker-Pattern', '07: Circuit Breaker Pattern', 'See markdown file', 7),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '08-Retries-and-Exponential-Backoff', '08: Retries and Exponential Backoff', 'See markdown file', 8),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '09-Timeouts-and-Deadlines', '09: Timeouts and Deadlines', 'See markdown file', 9),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '10-Distributed-Transactions', '10: Distributed Transactions', 'See markdown file', 10),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '11-Idempotency', '11: Idempotency', 'See markdown file', 11),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '12-Consensus-Algorithms', '12: Consensus Algorithms', 'See markdown file', 12),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '13-Distributed-Tracing', '13: Distributed Tracing', 'See markdown file', 13),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '14-Centralized-Logging', '14: Centralized Logging', 'See markdown file', 14),
(gen_random_uuid(), '88888888-8888-8888-8888-888888888888', '15-Metrics-and-Monitoring', '15: Metrics and Monitoring', 'See markdown file', 15)
ON CONFLICT (slug) DO NOTHING;
