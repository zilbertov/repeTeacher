UPDATE notifications
SET recipient_type = 'student'
WHERE description IN (
    'Новое сообщение от репетитора.',
    'Занятие перенесено репетитором.',
    'Занятие отменено репетитором.',
    'Занятие перенесено.',
    'Занятие отменено.'
);
