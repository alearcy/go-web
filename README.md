# GO WEB

## Local DB development

First, create user and database. From the terminal run:

`createuser [user_name]`

`createdb [db_name]`

then grant all privileges to the new user:

`psql [db_name]`

`GRANT ALL PRIVILEGES ON DATABASE [db_name] TO [user_name];`

exit

`\q`

To configure tables into the created database, execute `setup.sql` file:

`psql [db_name] -h [host] -a -f setup.sql`

## Styles development

Style is based on [Tailwindcss]('https://tailwindcss.com/') and can be run with:

`npm run start`

## Run Go

`go run main.go`
