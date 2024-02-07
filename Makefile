build :
	docker build -t go-app .

run :
	docker run -p 8080:8080 go-app

drop :
	sqlite3 forum.sqlite3 < migration/drop_all_table.sql 

migrate :
	sqlite3 forum.sqlite3 < migration/01.add_user_table.sql 
	sqlite3 forum.sqlite3 < migration/06.add_session_table.sql 
	sqlite3 forum.sqlite3 < migration/02.add_post_table.sql 
	sqlite3 forum.sqlite3 < migration/03.add_comment_table.sql 
	sqlite3 forum.sqlite3 < migration/04.add_post_categories_table.sql 
	sqlite3 forum.sqlite3 < migration/05.add_reaction_comments_table.sql 
	sqlite3 forum.sqlite3 < migration/07.add_reaction_post_table.sql 

create_categories : 
	sqlite3 forum.sqlite3 < migration/08.create_categories.sql

start : 
	go run ./cmd