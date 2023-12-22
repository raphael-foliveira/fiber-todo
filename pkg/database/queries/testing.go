package queries

const RecreateSchema = `
	DROP SCHEMA public CASCADE; 
	CREATE SCHEMA public; 
	SET search_path TO public;
`

const InsertTodoFixtures = `
	INSERT INTO todo 
		(title, description, completed) 
	VALUES 
		('test', 'test', false), 
		('test2', 'test2', false);
`

const ClearTodoTable = `
	DELETE FROM todo;
`
