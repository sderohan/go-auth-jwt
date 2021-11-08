# go-auth-jwt
## Here I have implemented the simple user login and logout system using JWT in golang.
## API endpoint information available below

### TECH STACK
- GO
   * fiber framework
   * jwt package
   * gorm for database manipulation
   * mail package for email validation
- MYSQL

### Requirements to run this programme 
- Make sure mysql database is running on the machine.
- Set the database connection information in `database/connection.go` file.
- Create the matching database name inside local database specified in `database/connection.go` file.
- Run the project by using the command `go run main.go`
- Alternate way to run using the `go build -o executable_file_name` and then run the `./executable_file_name`

### Following are the api endpoints available
1. Method `POST`
   * Endpoint `/api/register`
   * This endpoint allows to register the user with the provided information
   * Below is the accepted JSON format \
{ \
  "name": "valid username", \
  "email": "valid email id", \
  "password": "valid password" \
}

2. Method `POST`
   * Endpoint `/api/login`
   * This endpoint allows to login
   * This method returns the cookie in which JWT token is embeded

3. Method `GET`
   * Endpoint `/api/user`
   * This endpoint returns the currently logged in user information
   * Below is the data returned by this method \
{ \
  "ID": "Currently logged in user id"\
  "name": "Currently logged in username", \
  "email": "Currently logged in user email id", \
}

4. Method `Post`
   * Endpoint /api/logout
   * This endpoint clears the user session

`NOTE :` Below are the constraints for above fields
1. `name`
   * length of the name field should be greater than `four` and only uppercase and lowercase letter's are allowed
2. `email`
   * Entered email id should follow the normal valid email id format
3. `password`
   * length should be `greater than or equal to 8`
   * one uppercase letter
   * one lowercase letter
   * one digit
   * one special symbol from this list `! @ # $ % ^ &`
