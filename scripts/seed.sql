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

-- Seed user progress
INSERT INTO user_progress (user_id, entity_type, entity_id, status, completed_at) VALUES 
('11111111-1111-1111-1111-111111111111', 'course', '22222222-2222-2222-2222-222222222222', 'completed', NOW()),
('11111111-1111-1111-1111-111111111111', 'course', '33333333-3333-3333-3333-333333333333', 'in_progress', NULL)
ON CONFLICT DO NOTHING;
