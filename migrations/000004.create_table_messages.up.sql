CREATE TABLE IF NOT EXISTS public.messages (
    id BIGSERIAL PRIMARY KEY,
    chatId BIGINT NOT NULL,
    senderId INTEGER NOT NULL,
    content VARCHAR(255)
);
