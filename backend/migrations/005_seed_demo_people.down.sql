DELETE FROM lesson_files
WHERE lesson_id IN (
    SELECT l.id
    FROM lessons l
    JOIN students s ON s.id = l.student_id
    WHERE s.email = 'student.demo@example.com'
);

DELETE FROM lessons
WHERE student_id IN (
    SELECT id FROM students WHERE email = 'student.demo@example.com'
);

DELETE FROM student_subjects
WHERE student_id IN (
    SELECT id FROM students WHERE email = 'student.demo@example.com'
);

DELETE FROM students
WHERE email = 'student.demo@example.com';

DELETE FROM tutor_subjects
WHERE tutor_id IN (
    SELECT id FROM tutors WHERE email = 'demo.tutor@example.com'
);

DELETE FROM notification_settings
WHERE tutor_id IN (
    SELECT id FROM tutors WHERE email = 'demo.tutor@example.com'
);

DELETE FROM tutors
WHERE email = 'demo.tutor@example.com';
