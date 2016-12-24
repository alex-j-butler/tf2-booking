CREATE TABLE bookings (
    booking_id serial PRIMARY KEY,
    booker_id integer,
    server_name varchar (64) NOT NULL,
    booked_time timestamp,
    unbooked_time timestamp
);

CREATE INDEX server_name_idx ON bookings USING btree(server_name);

CREATE TABLE demos (
    demo_id serial PRIMARY KEY,
    booking_id integer,
    name varchar (64) NOT NULL,
    map_name varchar (32) NOT NULL,
    configuration varchar (32) NOT NULL,
    url varchar (128) NOT NULL,
    uploaded_time timestamp
);

CREATE TABLE users (
    user_id serial PRIMARY KEY,
    discord_id varchar (32),
    name varchar (32)
);

CREATE INDEX discord_id_idx ON users USING btree(discord_id);

CREATE TABLE demo_users (
    ref_id serial PRIMARY KEY,
    demo_id integer,
    user_id integer
);

ALTER TABLE bookings ADD FOREIGN KEY (booker_id) REFERENCES users (user_id);
ALTER TABLE demos ADD FOREIGN KEY (booking_id) REFERENCES bookings (booking_id);
ALTER TABLE demo_users ADD FOREIGN KEY (demo_id) REFERENCES demos;
ALTER TABLE demo_users ADD FOREIGN KEY (user_id) REFERENCES users;
