-- enum of roles
create type roles as ENUM('admin', 'client', 'worker');

create type o_status as ENUM('pending', 'processing', 'done');

CREATE table people (
	id serial primary key,
	name varchar(255),
	role roles,
	status varchar(255),
	email varchar(255),
	phone varchar(255),
	password text
);

CREATE table typework (
	id serial primary key,
	name varchar(255)
);

CREATE table orders (
	id serial primary key,
	model_name TEXT,
	warranty boolean,
	comment TEXT,
	client_id int,
	work_type int,
	worker_id int,
	order_status o_status,
	created_at TIMESTAMPTZ default NOW(),
	conf_time TIMESTAMPTZ,
	term int,
	summary int,
	foreign key (client_id) references people(id),
	foreign key (work_type) references typework(id)
);

CREATE table suggestions (
	id serial primary key,
	order_id int,
	worker_id int,
	summary int,
	term int,
	foreign key (worker_id) references people(id),
	foreign key (order_id) references orders(id)
);

create table account(
	id serial primary key,
	summary int,
	worker_id int,
	foreign key(worker_id)references people(id)
);

-- active, inactive
INSERT INTO typework(name) VALUES ('Другое'), ('Стиральная машина'), ('Телефон'), ('Игровая приставка'), ('Микроволновка'), ('Компьютер');

-- тестовые данные 
INSERT INTO people(name, role, email, phone, status, password) VALUES ('root', 'admin', 'artem.vecherinin@mail.ru', '89159045444', 'active', '$2a$14$5BjC3s/XPnifC.Dz7l30tOoc2Wfh1S3yYVhvvaNLaWmYW7vttWKQq'), ('Джон Ленон', 'client', 'arkel@gmail.com', '84133121212', 'active', '$2a$14$5BjC3s/XPnifC.Dz7l30tOoc2Wfh1S3yYVhvvaNLaWmYW7vttWKQq'), ('Роберт Полсон', 'worker', 'rober@gmail.com', '84231321212', 'active', '$2a$14$5BjC3s/XPnifC.Dz7l30tOoc2Wfh1S3yYVhvvaNLaWmYW7vttWKQq'), ('Герман Миллер', 'worker', 'rok@gmail.com', '84231321288', 'active', '$2a$14$5BjC3s/XPnifC.Dz7l30tOoc2Wfh1S3yYVhvvaNLaWmYW7vttWKQq');

-- уже полностью оформленный заказ.
INSERT INTO orders(model_name, warranty, comment, client_id, worker_id,  order_status, conf_time, term, summary, work_type) VALUES('Iphone', false, 'Broken screen', 2, 3, 'processing', now(), 3, 2000, 2);

-- Заказ который не был еще принят
INSERT INTO orders(model_name, warranty, comment, client_id, order_status, work_type, worker_id VALUES('Iphone', false, 'Broken screen', 2, 'pending', 2, null), ('Philips', true, 'broke boy', 2, 'pending', 2, 2), ('Sony', false, 'hueken', 2, 'pending', 2, 2), ('Kiki', true, 'broker max', 2, 'pending', 2, null);


INSERT INTO suggestions(order_id,  worker_id, summary, term) VALUES(2, 3, 1000, 20),(4, 2, 1000, 2);

ALTER SYSTEM SET summarize_wal = on;
ALTER SYSTEM SET shared_preload_libraries = 'pg_stat_statements';
ALTER SYSTEM SET pg_stat_statements.track = 'all';
ALTER SYSTEM SET track_activity_query_size = 2048;
SELECT pg_reload_conf();





-- INSERT INTO account(worker_id, summary) VALUES (1, 31000);
