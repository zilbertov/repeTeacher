ALTER TABLE notifications
DROP CONSTRAINT notifications_type_check;

ALTER TABLE notifications
ADD CONSTRAINT notifications_type_check
CHECK (type IN ('reschedule', 'cancel', 'message', 'new_request'));

ALTER TABLE tutors
DROP COLUMN notes;

