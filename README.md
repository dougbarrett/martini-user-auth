# Martini User Auth Example

This is an example Go Language application that uses the Martini framework, Gorm to handle database connections, and a few other middlewares.

### Key Features:

 - Encrypted passwords
 - Profile Creation Validation
 - Example on how to use layout.tmpl
 - Example on how to display data depending on if the user is logged in/out


### How to run:

1. Set up a new database, by default the config will look for the `demo` database
2. Create a new file named `config.toml` and use the contents of `config.toml.example` to start.  **Change the SqlConnection variable to your MySQL connection settings**
3. Run `go get` to grab all the dependencies
4. Run `go run *.go` or your favorite live-reload tool to test
5. If you use the default port, go to http://localhost:3003/ to test!
6. **The Test e-mail is: test@test.com, and password is `test123`**

### Questions? Comments?

Let me know if you find any bugs, or have any questions!  I will update this repo as needed!