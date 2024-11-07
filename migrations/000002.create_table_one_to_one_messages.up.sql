CREATE TABLE IF NOT EXISTS public.oneToOneMessages (
    messageId BIGINT PRIMARY KEY,
    senderId BIGINT NOT NULL,
    recipientId BIGINT NOT NULL,
    content VARCHAR(255),
    timestamp TIMESTAMP NOT NULL,
    CONSTRAINT fk_sender FOREIGN KEY (senderId) REFERENCES public.users(id),
    CONSTRAINT fk_recipient FOREIGN KEY (recipientId) REFERENCES public.users(id)
);
