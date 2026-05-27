ALTER TABLE notifications
ADD COLUMN recipient_type TEXT NOT NULL DEFAULT 'tutor'
CHECK (recipient_type IN ('tutor', 'student'));
