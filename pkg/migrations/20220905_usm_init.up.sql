create table tenants(
	id bigint primary key,
	code varchar(50) not null,
	label varchar(100) not null,
	status smallint default 0
);
create sequence organizations_id_seq as bigint increment by 1 minvalue 1 start with 1;

create table organizations (
	id bigint primary key default nextval('organizations_id_seq'),
	tenant_id bigint references tenants(id),
	code varchar(50) not null,
	label varchar(50) not null,
	type varchar(10) not null,
	status smallint default 0
);
create sequence sectors_id_seq as bigint increment by 1 minvalue 1 start with 1; 

create table sectors (
	id bigint primary key default nextval('sectors_id_seq'),
	tenant_id bigint not null not null references tenants(id),
	org_id bigint not null not null references organizations(id),
	code varchar(50) not null,
	label varchar(50) not null,
	parent_id bigint references sectors(id),
	has_parent boolean default false,
	depth smallint default 0,
	status smallint default 0
);
create sequence users_id_seq as bigint increment by 1 minvalue 1 start with 1;

create table users(
	id bigint primary key default nextval('users_id_seq'),
	tenant_id bigint not null references tenants(id),
	org_id bigint not null references organizations(id),
	external_id varchar(50) not null,
	last_name varchar(50) not null,
	first_name varchar(50) not null,
	middle_name varchar(50),
	login varchar(50) not null,
	email varchar(50) not null,
	status smallint default 0
);

insert into tenants (id,code,"label") values (1,'lxc','LXConnect');
