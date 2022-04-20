package sqlite

const SQL_SCHEMA = `    

DROP TABLE IF EXISTS config;
CREATE TABLE IF NOT EXISTS config (
	key VARCHAR NOT NULL, 
	value VARCHAR NOT NULL, 
	PRIMARY KEY (key)
);

CREATE TABLE workflow_definitions (
	id VARCHAR NOT NULL, 
	created DATETIME NOT NULL, 
	last_modified DATETIME NOT NULL, 
	version INTEGER NOT NULL, 
	yaml VARCHAR NOT NULL, 
	is_deleted BOOLEAN NOT NULL, 
	is_enabled BOOLEAN NOT NULL, 
	deleted_date DATETIME, 
	PRIMARY KEY (id), 
	CHECK (is_deleted IN (0, 1)), 
	CHECK (is_enabled IN (0, 1))
);

CREATE TABLE workflow_history (
	id VARCHAR NOT NULL, 
	version INTEGER NOT NULL, 
	created DATETIME NOT NULL, 
	reason VARCHAR NOT NULL, 
	yaml VARCHAR NOT NULL, 
	PRIMARY KEY (id, version)
);

CREATE TABLE workflow_instances (
	id VARCHAR NOT NULL, 
	created DATETIME NOT NULL, 
	last_modified DATETIME NOT NULL, 
	finished DATETIME, 
	is_active BOOLEAN NOT NULL, 
	outcome VARCHAR NOT NULL, 
	state VARCHAR NOT NULL, 
	yaml VARCHAR NOT NULL, 
	workflow_id VARCHAR, 
	shared_state_network_id VARCHAR, 
	shared_state_volume_id VARCHAR, 
	shared_state_container_id VARCHAR, 
	PRIMARY KEY (id), 
	CHECK (is_active IN (0, 1)), 
	FOREIGN KEY(workflow_id) REFERENCES workflow_definitions (id)
);

`
