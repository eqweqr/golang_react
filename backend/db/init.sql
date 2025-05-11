-- enum of roles
create type roles as ENUM('admin', 'client', 'worker');

create type o_status as ENUM('pending', 'processing', 'done', 'return', 'approve');

CREATE table people (
	id serial primary key,
	name varchar(255),
	role roles,
	email varchar(255),
	phone varchar(255),
	password text
);

CREATE table orders (
	id serial primary key,
	model_name varchar(255),
	warranty boolean,
	comment text,
	summary int,
	client_id int,
	order_status o_status,
	worker_id int,
	foreign key (client_id) references people(id),
	foreign key (worker_id) references people(id) 
);

create table account(
	id serial primary key,
	worker_id int,
	foreign key(worker_id) references people(id),
	summary int
);


-- тестовые данные 
INSERT INTO people(name, role, email, phone) VALUES ('root', 'admin', 'artem.vecherinin@mail.ru', '89159045444'), ('Джон Ленон', 'client', 'arkel@gmail.com', '84133121212'), ('Роберт Полсон', 'worker', 'rober@gmail.com', '84231321212');

INSERT INTO orders(model_name, warranty, comment, summary, client_id, order_status, worker_id) VALUES('Model1', false, 'This my second q', 2000, 3, 'pending', 2);



INSERT INTO account(worker_id, summary) VALUES (1, 31000);