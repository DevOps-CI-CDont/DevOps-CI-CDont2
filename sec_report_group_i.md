# Security report group i

## risk identification

### assets

* uptime of your web app.
  * in reality downtime would cost money, but for this project it's just a matter of pride.
* user data - emails, passwords (hopefully hashed), etc.

### sources

* SQL Injection if you don't sanitize your inputs or use an ORM that does this for you.
* unencrypted traffic (HTTP) between client and server.
* virtual machine intrusion.

### scenarios

* SQL Injection
  * attacker can read all data from the database
  * attacker can modify data in the database
  * attacker can delete data from the database
  * attacker can create new users (without using frontend/api)
  * attacker can create new posts (without using frontend/api)
* Virtual Machine intrusion
  * if root access is achieved inside VM, attacker has full control over the webapp.

* Unencrypted traffic (HTTP) between client and server
  * attacker can read all data sent between client and server
  * attacker can modify data sent between client and server

## risk analysis

### likelihood

* SQL Injection
  * low
    * SQL Injection is a very common attack vector, but it's not very likely that the attacker will be able to find a vulnerability in the webapp.
* DOS (Denial of Service)
  * low - medium


### impact

* SQL Injection
  * high
    * previous section about SQL Injection scenario explains the impact of this attack.

## pentesting

`wmap` found endpoints:
[+] Found http://157.245.16.6:8000/docs 200
[+] Found http://157.245.16.6:8000/login 200
[+] Found http://157.245.16.6:8000/logout 302
[+] Found http://157.245.16.6:8000/public 200
[+] Found http://157.245.16.6:8000/register 200
