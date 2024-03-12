CREATE TABLE BANKS (
    ID SERIAL PRIMARY KEY,
    NAME VARCHAR(15) NOT NULL,
    ACCOUNT_NAME VARCHAR(15) NOT NULL,
    ACCOUNT_NUMBER VARCHAR(15) NOT NULL,
    USER_ID INT NOT NULL,
    CREATED_AT BIGINT NOT NULL,
    UPDATED_AT BIGINT NOT NULL,
    CONSTRAINT fk_banks_user FOREIGN KEY(USER_ID) REFERENCES USERS(id)
);