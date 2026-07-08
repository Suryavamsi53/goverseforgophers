-- Ensure the 'postgresql' course exists
INSERT INTO courses (id, slug, title, description, difficulty) 
VALUES ('c0000000-0000-0000-0000-000000000002', 'postgresql', 'PostgreSQL Mastery', 'Master PostgreSQL from basics to advanced real-world production scaling.', 'expert')
ON CONFLICT (id) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description;

-- Insert the massive single roadmap lesson at order_index 0 so it's the first thing they see
INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES
('30000000-0000-0000-0000-000000000000', 'c0000000-0000-0000-0000-000000000002', '000-roadmap', 'PostgreSQL Complete Roadmap & Scenarios', 'See markdown file', 0)
ON CONFLICT (id) DO UPDATE SET slug = EXCLUDED.slug, title = EXCLUDED.title, content = EXCLUDED.content, order_index = EXCLUDED.order_index;
