Command for installing database on local:
1. Install PostgreSQL via brew
`brew update`
`brew install postgresql@14`
2. Start the PostgreSQL server 
`brew services start postgresql@14`
3. Verify installation & add to PATH (if needed)
`   psql --version`
If not found, ensure PostgreSQL’s bin folder is in your PATH. Example for version 14:
`echo 'export PATH="/usr/local/opt/postgresql@14/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc`
4. Log in using psql
   A. Connect to the default postgres server:
   `psql postgres`
   or specify username (typically your macOS login name):
`psql -U postgres -d postgres`
To prompt for password:
`psql -U your_username -d your_db -W`
5. Create your own user/role and database
   Once inside psql
`   CREATE ROLE myuser WITH LOGIN PASSWORD 'securepass';
   ALTER ROLE myuser CREATEDB;       -- allows new DB creation`
   Quit psql with \q and log in as the new user:
`psql -U myuser -d postgres`
Then create your own database:
`CREATE DATABASE task_manager;
GRANT ALL PRIVILEGES ON DATABASE task_manager TO myuser;`