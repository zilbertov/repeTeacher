INSERT INTO students (tutor_id, name, email, phone, exam_type, status, notes)
SELECT 1, 'Тестовый Ученик', 'student.demo@example.com', '89000000000', 'ЕГЭ', 'active', 'Демо-ученик для показа роли ученика.'
WHERE NOT EXISTS (
    SELECT 1 FROM students WHERE tutor_id = 1 AND email = 'student.demo@example.com'
);

INSERT INTO student_subjects (student_id, subject)
SELECT s.id, 'Математика'
FROM students s
WHERE s.tutor_id = 1
  AND s.email = 'student.demo@example.com'
  AND NOT EXISTS (
      SELECT 1 FROM student_subjects ss
      WHERE ss.student_id = s.id AND ss.subject = 'Математика'
  );

INSERT INTO lessons (tutor_id, student_id, subject, exam_type, lesson_date, start_time, duration_minutes, format, has_homework, price)
SELECT 1, s.id, 'Математика', 'ЕГЭ', '2026-04-02', '10:00', 60, 'online', true, 900
FROM students s
WHERE s.tutor_id = 1
  AND s.email = 'student.demo@example.com'
  AND NOT EXISTS (
      SELECT 1 FROM lessons l
      WHERE l.tutor_id = 1
        AND l.student_id = s.id
        AND l.lesson_date = '2026-04-02'
        AND l.start_time = '10:00'
  );

SELECT setval(pg_get_serial_sequence('tutors', 'id'), COALESCE((SELECT MAX(id) FROM tutors), 1), true);

INSERT INTO tutors (name, email, phone, password_hash)
SELECT 'Анна Смирнова', 'demo.tutor@example.com', '89191234567', 'demo-password-hash'
WHERE NOT EXISTS (
    SELECT 1 FROM tutors WHERE email = 'demo.tutor@example.com'
);

INSERT INTO notification_settings (tutor_id)
SELECT t.id
FROM tutors t
WHERE t.email = 'demo.tutor@example.com'
  AND NOT EXISTS (
      SELECT 1 FROM notification_settings ns WHERE ns.tutor_id = t.id
  );

INSERT INTO tutor_subjects (tutor_id, subject)
SELECT t.id, 'Математика'
FROM tutors t
WHERE t.email = 'demo.tutor@example.com'
  AND NOT EXISTS (
      SELECT 1 FROM tutor_subjects ts
      WHERE ts.tutor_id = t.id AND ts.subject = 'Математика'
  );
