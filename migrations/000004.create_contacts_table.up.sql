CREATE TABLE IF NOT EXISTS public.contacts (
    userId BIGINT NOT NULL,
    contactId BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    isDeleted BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (userId, contactId),
    CONSTRAINT fk_contact FOREIGN KEY (contactId) REFERENCES public.users(id),
    CONSTRAINT fk_user FOREIGN KEY (userId) REFERENCES public.users(id)
);
