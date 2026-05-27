CREATE TABLE tutor_subjects (
    id BIGSERIAL PRIMARY KEY,
    tutor_id BIGINT NOT NULL REFERENCES tutors(id) ON DELETE CASCADE,
    subject TEXT NOT NULL
);

INSERT INTO tutor_subjects (tutor_id, subject)
VALUES
    (1, 'Русский язык'),
    (1, 'Математика');
