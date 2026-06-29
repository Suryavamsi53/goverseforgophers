-- Seed an initial admin user
INSERT INTO users (id, username, email, password_hash, role) 
VALUES ('11111111-1111-1111-1111-111111111111', 'gopher_master', 'master@goverse.dev', 'hash', 'admin')
ON CONFLICT DO NOTHING;

-- Seed user profile
INSERT INTO user_profiles (user_id, avatar_url, bio, github_handle, daily_streak, total_score)
VALUES ('11111111-1111-1111-1111-111111111111', 'https://github.com/golang.png', 'Senior Gopher', 'golang', 14, 1337)
ON CONFLICT DO NOTHING;

-- Seed some courses
INSERT INTO courses (id, slug, title, description, difficulty) VALUES 
('22222222-2222-2222-2222-222222222222', 'go-fundamentals', 'Go Fundamentals', 'Learn the basics of Go.', 'beginner'),
('33333333-3333-3333-3333-333333333333', 'go-concurrency', 'Mastering Concurrency', 'Deep dive into Goroutines and Channels.', 'advanced')
ON CONFLICT DO NOTHING;

-- Seed some lessons
INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
('44444444-4444-4444-4444-444444444444', '22222222-2222-2222-2222-222222222222', 'basics', 'Go Basics', 'Learn variables, types, and control structures.', 1),
('55555555-5555-5555-5555-555555555555', '22222222-2222-2222-2222-222222222222', 'structs', 'Structs and Methods', 'Learn composite types and methods.', 2),
('66666666-6666-6666-6666-666666666666', '33333333-3333-3333-3333-333333333333', 'concurrency', 'Go Concurrency', 'Master Goroutines, Channels, and the Select statement.', 1)
ON CONFLICT DO NOTHING;

-- Seed user progress
INSERT INTO user_progress (user_id, entity_type, entity_id, status, completed_at) VALUES 
('11111111-1111-1111-1111-111111111111', 'course', '22222222-2222-2222-2222-222222222222', 'completed', NOW()),
('11111111-1111-1111-1111-111111111111', 'course', '33333333-3333-3333-3333-333333333333', 'in_progress', NULL),
('11111111-1111-1111-1111-111111111111', 'lesson', '44444444-4444-4444-4444-444444444444', 'completed', NOW())
ON CONFLICT DO NOTHING;
