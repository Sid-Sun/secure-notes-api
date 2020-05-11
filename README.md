# PASSWORDLESS SECURE NOTES API

An API written in Go to demonstrate CRUD (Create Read Update Delete) on encrypted notes anonymous without storing password in any manner - As an example to using cryptography to facilitate authentication

A proof-of-concept to my Medium article on [Doing secure password authentication without storing passwords](https://medium.com/@sidharth.soni525/doing-secure-password-authentication-without-storing-passwords-part-1-7b6024843763) - the main difference being here we do CRUD on notes instead of encryption and decryption on asymmetric private keys

# Usage:

## To Be Greeted

We all like greetings - This API likes to greet, but it is not pushy about it, it only offers you its heartfelt time-appropriate greetings if you ask for them. 
Request Type: `HTTP GET`

URL:
```
http://$host:$port/
```



## To Set A New Note

Request Type: `HTTP POST`

BODY:
```json
{
    "id": "mySecretNote",
    "pass": "&now:we@pluto",
    "note": "I am a bufferfly, flying through the sky"
}
```

URL:
```
http://$host:$port/set
```


Upon success, the Original ID is returned. **If an ID is not supplied or a note already exists with the supplied id**, A randomly generated 8 character ID is returned insted.

Response on success:
```json
{
    "ID": "mySecretNote"
}
```


A request may fail due to any one (or all) of the following reasons:

  * The note in request is empty
  * The password in request is empty

If the Request fails, the response looks like:
```json
{
    "ID": ""
}
```



## To Get A Note

Request Type: `HTTP GET`

BODY:
```json
{
    "id": "mySecretNote",
    "pass": "&now:we@pluto"
}
```

URL:
```
http://$host:$port/get
```


Upon success, the Original ID is returned along with the note.

Response on success:
```json
{
    "ID": "mySecretNote",
    "Note": "I am a bufferfly, flying through the sky"
}
```


A request may fail due to any one of the following reasons:

  1. The ID in the request is empty
  2. The pass in request is empty
  3. The supplied ID is incorrect
  4. The supplied pass is incorrect

If the Request fails due to 1 or 2, the response looks like:
```json
{
    "ID": "",
    "Note": ""
}
```

If the Request fails due to 3 or 4, the response looks like:
```json
{
    "ID": "mySecretNote",
    "Note": ""
}
```



## To Delete A Note

Request Type: `HTTP DELETE`

BODY:
```json
{
    "id": "mySecretNote",
    "pass": "&now:we@pluto"
}
```

URL:
```
http://$host:$port/delete
```


Upon success, the Original ID is returned.

Response on success:
```json
{
    "ID": "mySecretNote"
}
```


A request may fail due to any one of the following reasons:

  * The ID in the request is empty
  * The pass in request is empty
  * The supplied ID is incorrect
  * The supplied pass is incorrect

If the Request fails, the response looks like:
```json
{
    "ID": ""
}
```



## To Update A Note

Request Type: `HTTP PUT`

BODY:
```json
{
    "id": "mySecretNote",
    "pass": "&now:we@pluto",
    "note": "I am a bufferfly, flying through the sky on Mars"
}
```

Optionally to change the password for the new note, newpass may also be defined, like so:

BODY:
```json
{
    "id": "mySecretNote",
    "pass": "&now:we@pluto",
    "newpass": "&now:we@moon",
    "note": "I am a bufferfly, flying through the sky on Mars"
}
```

URL:
```
http://$host:$port/update/note
```


Upon success, the Original ID is returned.

Response on success:
```json
{
    "ID": "mySecretNote"
}
```


A request may fail due to any one of the following reasons:

  * The ID, pass or note in the request are empty
  * The supplied ID or pass are incorrect

If the Request fails, the response looks like:
```json
{
    "ID": ""
}
```



## To Update A Note's Pass

Request Type: `HTTP PATCH`

BODY:
```json
{
    "id": "mySecretNote",
    "pass": "&now:we@pluto",
    "newpass": "&now:we@moon"
}
```

URL:
```
http://$host:$port/update/pass
```


Upon success, the Original ID is returned.

Response on success:
```json
{
    "ID": "mySecretNote"
}
```


A request may fail due to any one of the following reasons:

  * The ID, pass or newpass fields are request are empty
  * The newpass and pass fields are equal
  * The supplied ID or pass are incorrect

If the Request fails, the response looks like:
```json
{
    "ID": ""
}
```

## Cheers!