CREATE TABLE tutors (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    tutor_id BIGINT NOT NULL REFERENCES tutors(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    exam_type TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('request', 'active', 'archived')),
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE student_subjects (
    id BIGSERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    subject TEXT NOT NULL
);

CREATE TABLE notification_settings (
    tutor_id BIGINT PRIMARY KEY REFERENCES tutors(id) ON DELETE CASCADE,
    push_enabled BOOLEAN NOT NULL DEFAULT true,
    telegram_enabled BOOLEAN NOT NULL DEFAULT true,
    sound_enabled BOOLEAN NOT NULL DEFAULT true,
    lesson_reminders_enabled BOOLEAN NOT NULL DEFAULT false
);

INSERT INTO tutors (id, name, email, phone, password_hash)
VALUES (1, 'Вадим Зильбертов', 'v4bem@ya.ru', '89198318673', 'demo-password-hash');

INSERT INTO notification_settings (tutor_id)
VALUES (1);
