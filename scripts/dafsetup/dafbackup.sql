--
-- Copyright (c) 2019 by Delphix. All rights reserved.
--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.8
-- Dumped by pg_dump version 9.6.8

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: :user; Type: USER;
--
CREATE USER :user WITH
	LOGIN
	SUPERUSER
	CREATEDB
	CREATEROLE
	INHERIT
	REPLICATION
	CONNECTION LIMIT -1
	PASSWORD :'pass';

--
-- Name: dafdb; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE dafdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';

ALTER DATABASE dafdb OWNER TO :user;

\connect dafdb

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: billings; Type: TABLE; Schema: public; Owner: :user
--

-- CREATE TABLE public.billings (
--     id integer NOT NULL,
--     patient_id integer NOT NULL,
--     ccnum character varying(32) NOT NULL,
--     cctype character varying(55) NOT NULL,
--     ccexpmonth integer NOT NULL,
--     ccexpyear integer NOT NULL,
--     address1 character varying(255) NOT NULL,
--     address2 character varying(255),
--     city character varying(255) NOT NULL,
--     state character varying(4) NOT NULL,
--     zip character varying(12) NOT NULL
-- );

CREATE TABLE public.billings (
    id integer NOT NULL,
    patient_id integer NOT NULL,
    ccnum character varying(32) NOT NULL,
    cctype character varying(55) NOT NULL,
    ccexpmonth integer NOT NULL,
    ccexpyear integer NOT NULL,
    address1 character varying(255) NOT NULL,
    address2 character varying(255),
    city character varying(255) NOT NULL,
    state character varying(4) NOT NULL,
    zip character varying(12) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.billings OWNER TO :user;

--
-- Name: billings_id_seq; Type: SEQUENCE; Schema: public; Owner: :user
--

CREATE SEQUENCE public.billings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.billings_id_seq OWNER TO :user;

--
-- Name: billings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: :user
--

ALTER SEQUENCE public.billings_id_seq OWNED BY public.billings.id;


--
-- Name: databasechangelog; Type: TABLE; Schema: public; Owner: :user
--

CREATE TABLE public.databasechangelog (
    id character varying(255) NOT NULL,
    author character varying(255) NOT NULL,
    filename character varying(255) NOT NULL,
    dateexecuted timestamp without time zone NOT NULL,
    orderexecuted integer NOT NULL,
    exectype character varying(10) NOT NULL,
    md5sum character varying(35),
    description character varying(255),
    comments character varying(255),
    tag character varying(255),
    liquibase character varying(20),
    contexts character varying(255),
    labels character varying(255),
    deployment_id character varying(10)
);


ALTER TABLE public.databasechangelog OWNER TO :user;

--
-- Name: databasechangeloglock; Type: TABLE; Schema: public; Owner: :user
--

CREATE TABLE public.databasechangeloglock (
    id integer NOT NULL,
    locked boolean NOT NULL,
    lockgranted timestamp without time zone,
    lockedby character varying(255)
);


ALTER TABLE public.databasechangeloglock OWNER TO :user;

--
-- Name: patients; Type: TABLE; Schema: public; Owner: :user
--

-- CREATE TABLE public.patients (
--     id integer NOT NULL,
--     firstname character varying(255) NOT NULL,
--     middlename character varying(255) NOT NULL,
--     lastname character varying(255) NOT NULL,
--     ssn character varying(12) NOT NULL,
--     dobyear smallint NOT NULL,
--     dobmonth smallint NOT NULL,
--     dobday smallint NOT NULL,
--     address1 character varying(255) NOT NULL,
--     address2 character varying(255) NOT NULL,
--     city character varying(255) NOT NULL,
--     state character varying(4) NOT NULL,
--     zip character varying(12) NOT NULL
-- );

CREATE TABLE public.patients (
    id integer NOT NULL,
    firstname character varying(255) NOT NULL,
    middlename character varying(255) NOT NULL,
    lastname character varying(255) NOT NULL,
    ssn character varying(12) NOT NULL,
    dobyear smallint NOT NULL,
    dobmonth smallint NOT NULL,
    dobday smallint NOT NULL,
    address1 character varying(255) NOT NULL,
    address2 character varying(255) NOT NULL,
    city character varying(255) NOT NULL,
    state character varying(4) NOT NULL,
    zip character varying(12) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

ALTER TABLE public.patients OWNER TO :user;

--
-- Name: patients_id_seq; Type: SEQUENCE; Schema: public; Owner: :user
--

CREATE SEQUENCE public.patients_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.patients_id_seq OWNER TO :user;

--
-- Name: patients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: :user
--

ALTER SEQUENCE public.patients_id_seq OWNED BY public.patients.id;


--
-- Name: payments; Type: TABLE; Schema: public; Owner: :user
--

-- CREATE TABLE public.payments (
--     id integer NOT NULL,
--     patient_id integer NOT NULL,
--     amount integer NOT NULL,
--     authcode character varying(36) NOT NULL,
--     currency character varying(4) NOT NULL,
--     captured boolean DEFAULT false NOT NULL,
--     type character varying(55) NOT NULL
-- );

CREATE TABLE public.payments (
    id integer NOT NULL,
    patient_id integer NOT NULL,
    amount integer NOT NULL,
    authcode character varying(36) NOT NULL,
    currency character varying(4) NOT NULL,
    captured boolean DEFAULT false NOT NULL,
    type character varying(55) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

ALTER TABLE public.payments OWNER TO :user;

--
-- Name: payments_id_seq; Type: SEQUENCE; Schema: public; Owner: :user
--

CREATE SEQUENCE public.payments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.payments_id_seq OWNER TO :user;

--
-- Name: payments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: :user
--

ALTER SEQUENCE public.payments_id_seq OWNED BY public.payments.id;


--
-- Name: records; Type: TABLE; Schema: public; Owner: :user
--

-- CREATE TABLE public.records (
--     id integer NOT NULL,
--     patient_id integer NOT NULL,
--     type character varying(255) NOT NULL
-- );

CREATE TABLE public.records (
    id integer NOT NULL,
    patient_id integer NOT NULL,
    type character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

ALTER TABLE public.records OWNER TO :user;

--
-- Name: records_id_seq; Type: SEQUENCE; Schema: public; Owner: :user
--

CREATE SEQUENCE public.records_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.records_id_seq OWNER TO :user;

--
-- Name: records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: :user
--

ALTER SEQUENCE public.records_id_seq OWNED BY public.records.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: :user
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(255) NOT NULL,
    firstname character varying(255) NOT NULL,
    lastname character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    password character varying(255) NOT NULL
);

ALTER TABLE public.users OWNER TO :user;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: :user
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO :user;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: :user
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: billings id; Type: DEFAULT; Schema: public; Owner: :user
--

ALTER TABLE ONLY public.billings ALTER COLUMN id SET DEFAULT nextval('public.billings_id_seq'::regclass);


--
-- Name: patients id; Type: DEFAULT; Schema: public; Owner: :user
--

ALTER TABLE ONLY public.patients ALTER COLUMN id SET DEFAULT nextval('public.patients_id_seq'::regclass);


--
-- Name: payments id; Type: DEFAULT; Schema: public; Owner: :user
--

ALTER TABLE ONLY public.payments ALTER COLUMN id SET DEFAULT nextval('public.payments_id_seq'::regclass);


--
-- Name: records id; Type: DEFAULT; Schema: public; Owner: :user
--

ALTER TABLE ONLY public.records ALTER COLUMN id SET DEFAULT nextval('public.records_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: :user
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: billings; Type: TABLE DATA; Schema: public; Owner: :user
--

\COPY public.billings (id, patient_id, ccnum, cctype, ccexpmonth, ccexpyear, address1, address2, city, state, zip, created_at, updated_at) FROM '/tmp/dafsetup/billings.csv'DELIMITER ',' CSV HEADER;;


--
-- Name: billings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: :user
--

SELECT pg_catalog.setval('public.billings_id_seq', 1, false);


--
-- Data for Name: patients; Type: TABLE DATA; Schema: public; Owner: :user
--

\COPY public.patients (id, firstname, middlename, lastname, ssn, dobyear, dobmonth, dobday, address1, address2, city, state, zip, created_at, updated_at) FROM '/tmp/dafsetup/patients.csv'DELIMITER ',' CSV HEADER;


--
-- Name: patients_id_seq; Type: SEQUENCE SET; Schema: public; Owner: :user
--

SELECT pg_catalog.setval('public.patients_id_seq', 1, false);


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: :user
--

\COPY public.payments (id, patient_id, amount, authcode, currency, captured, type, created_at, updated_at) FROM '/tmp/dafsetup/payments.csv'DELIMITER ',' CSV HEADER;


--
-- Name: payments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: :user
--

SELECT pg_catalog.setval('public.payments_id_seq', 1, false);


--
-- Data for Name: records; Type: TABLE DATA; Schema: public; Owner: :user
--

\COPY public.records (id, patient_id, type, created_at, updated_at) FROM '/tmp/dafsetup/records.csv'DELIMITER ',' CSV HEADER;;


--
-- Name: records_id_seq; Type: SEQUENCE SET; Schema: public; Owner: :user
--

SELECT pg_catalog.setval('public.records_id_seq', 1, false);


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: :user
--

COPY public.users (id, username, firstname, lastname, password) FROM stdin;
\.


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: :user
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- Name: billings billings_pkey; Type: CONSTRAINT; Schema: public; Owner: :user
--

ALTER TABLE ONLY public.billings
    ADD CONSTRAINT billings_pkey PRIMARY KEY (id);


--
-- Name: patients patients_pkey; Type: CONSTRAINT; Schema: public; Owner: :user
--

ALTER TABLE ONLY public.patients
    ADD CONSTRAINT patients_pkey PRIMARY KEY (id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: :user
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: records records_pkey; Type: CONSTRAINT; Schema: public; Owner: :user
--

ALTER TABLE ONLY public.records
    ADD CONSTRAINT records_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: :user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

