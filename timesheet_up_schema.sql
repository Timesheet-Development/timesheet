--Create schema for timesheet application.

CREATE TABLE IF NOT EXISTS users (
    id uuid not null primary key,
    login_name varchar(60) not null, 
    password varchar(60) not null,
    Name varchar(60) not null,
    address varchar(100) not null,
    department varchar(20) not null,
    security_no int4 not null,
    dob date not null,
    city varchar(60) null,
    state varchar(60) null,
    job_title varchar(60) not null,
    is_perm boolean not null,
    gender varchar(10) not null,
    passport varchar(60) not null,
    reporting_mngr uuid null,
    created_at date default now(),
    updated_at date default now()
);