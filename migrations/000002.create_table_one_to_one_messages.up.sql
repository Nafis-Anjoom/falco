CREATE TABLE IF NOT EXISTS public.oneToOneMessages (
    id BIGSERIAL PRIMARY KEY,
    senderId INTEGER NOT NULL,
    receiverId INTEGER NOT NULL,
    content VARCHAR(255),
    timestamp TIMESTAMP NOT NULL,
    CONSTRAINT fk_sender FOREIGN KEY (senderId) REFERENCES public.users(id),
    CONSTRAINT fk_receiver FOREIGN KEY (receiverId) REFERENCES public.users(id)
);
