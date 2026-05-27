CREATE TABLE lessons (
    id BIGSERIAL PRIMARY KEY,
    tutor_id BIGINT NOT NULL REFERENCES tutors(id) ON DELETE CASCADE,
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    subject TEXT NOT NULL,
    exam_type TEXT NOT NULL,
    lesson_date DATE NOT NULL,
    start_time TIME NOT NULL,
    duration_minutes INTEGER NOT NULL,
    format TEXT NOT NULL CHECK (format IN ('online', 'offline')),
    has_homework BOOLEAN NOT NULL DEFAULT false,
    price INTEGER NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('planned', 'cancelled', 'done')) DEFAULT 'planned',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE lesson_files (
    id BIGSERIAL PRIMARY KEY,
    lesson_id BIGINT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    file_type TEXT NOT NULL CHECK (file_type IN ('material', 'homework')),
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
