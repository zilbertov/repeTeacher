ALTER TABLE tutors
ADD COLUMN notes TEXT NOT NULL DEFAULT '';

ALTER TABLE notifications
DROP CONSTRAINT notifications_type_check;

ALTER TABLE notifications
ADD CONSTRAINT notifications_type_check
CHECK (type IN ('reschedule', 'cancel', 'message', 'new_request', 'lesson'));

