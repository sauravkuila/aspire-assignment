--drop tables for fresh schema
DROP TYPE IF EXISTS UserTypes;
DROP TYPE IF EXISTS LoanStatus;
DROP TYPE IF EXISTS LoanTransactionStatus;
DROP TABLE IF EXISTS user_detail;
DROP TABLE IF EXISTS loan;
DROP TABLE IF EXISTS installment;

--create types
CREATE TYPE UserTypes AS ENUM('USER','APPROVER');
CREATE TYPE LoanStatus AS ENUM('PENDING','APPROVED','INFORCE','REJECTED','CANCELLED','SETTLED');
CREATE TYPE LoanTransactionStatus AS ENUM('PENDING','PAID','CANCELLED');

-- create a function for timestamp
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

--create tables
CREATE TABLE user_detail(
   id serial,
   user_type UserTypes not null,
   email text not null, 
   mobile text not null,
   monthly_salary float DEFAULT 0.0,
   acc_bal float DEFAULT 0.0,
   created_at timestamp DEFAULT CURRENT_TIMESTAMP,
   updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY(id)
);

CREATE TABLE loan(
    id serial,
    user_id int,
    amount float not null,
    installments int not null,
    status LoanStatus not null,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    CONSTRAINT fk_userid
   		FOREIGN KEY(user_id) 
		REFERENCES user_detail(id)
);

CREATE TABLE installment(
    id serial,
    loan_id int,
    amount_due float not null,
    amount_paid float default 0,
    status LoanTransactionStatus not null,
    installment_num int not null,
    due_date timestamp not null,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    CONSTRAINT fk_loanid
   		FOREIGN KEY(loan_id) 
		REFERENCES loan(id)
);

-- create a trigger for timestamp
CREATE TRIGGER set_timestamp
AFTER UPDATE ON user_detail
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
AFTER UPDATE ON loan
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
AFTER UPDATE ON installment
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
