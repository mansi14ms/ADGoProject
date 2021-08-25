The project contains a list of APIs in Golang that interacts through a Non-SQL Database, MongoDb(Atlas-CloudDB)
### APIs

- Display all documents in database

  `GET : localhost/10000/all`
- Create New User
(Adds new user data only if the user's info is not already present)

  `POST : localhost/10000/addUserDB`
- Display the document with given id

  `GET : localhost/10000/getUser/{id}`
- Delete the document with given id

  `DELETE : localhost/10000/deleteUserFromId/{id}`
- Update the document
(Updates only if the id is present)

  `POST : localhost/10000/updateUserDataDB/{id}`
- Custom  update of the document
(Updates only if the id is present)

  `PUT : localhost/10000/updateUserData/{id}`


##### The entire documentation of the APIs can be found at : [GoAPI](https://documenter.getpostman.com/view/12416782/TzzGFt9u/)
