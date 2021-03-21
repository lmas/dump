MANAGE=python djangoforum/manage.py

runserver:
	python djangoforum/manage.py runserver

reset:
	rm -f djangoforum/db.sqlite3
	$(MANAGE) syncdb --noinput
	$(MANAGE) loaddata djangoforum/fixture.json
