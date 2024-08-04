# aspire-assignment
The repository provides an application that allows authenticated users to go through a loan application.
The assignment covers
* users applying for loan
* admins approving/rejecting loan
* users can submit repayments
There are a few assumptions and borders drawn around the applicaiotn in the interest of time and scoping. 

## Mandatory Features (as requested by assignment provider)
* Customer loans need to approved by Admin
* A loan must have an amount and loan term
* The loan and scheduled payments will have state `PENDING`
* On Admin approving the loan, the loan state changes from `PENDING` to `APPROVED`
* Policy check for cutomers to be able to view their own loan only
* Customer can add a repayment greater or equal to the scheduled repayment
* On repayment, the scheduled repayment state changes to `PAID`
* If all the scheduled repayments connected to a loan are `PAID`, automatically the loan also become `PAID`

## Additional Fetures Added
* A JWT based auth management added for customers/admins to signup and login
* Customers are allowed to not only apply for a loan but also modify the tenure, amount and even cancel the loan.
* No installments are generated for a loan unless approved by Admin
* Admins can not only `APPROVE` but also `REJECT` a loan
* Admins can list all loans which are in `PENDING` state to decide which takes priority of approval/rejection
* All loans/installments are tracked when they were created/approved/paid
* Customer can repay an amount equal or more than scheduled payment and the upcoming scheduled payments are adjusted equally.
* Customer can close the loan by making greater payments vs the scheduled payment amount

## Assumptions
* All loans will be assumed to have weekly payment frequency
* All loans provided are zero-interest loans
* Admins cannot apply for loan using the applicaiton

---

## Dependecies
the applicaiton uses postgres connection to manage tables. Use of a psql Docker image or local installation is required

### Docker usage
run a docer image for postgres
```
# docker run --name psql -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:12-alpine
# docker exec -it psql bash
# psql -h localhost -U postgres
# CREATE DATABASE aspire;
# \l
```
this should list the newly created database ```aspire```
connect to database created and execute the ```create.sql``` available in ```pkg/db/sql/create.sql``` path
```
# \c aspire;
# <run contents of> create.sql
# \d
```

### Local Server / Remote Server
connect to the local / remote server using a client like TablePlus, DBeaver, PostgreSQL.
execute the following in the SQL editor
```
CREATE DATABASE aspire;
<run contents of> create.sql
```

---

## Prerequisites
* run the db schema available in ```pkg/db/sql/create.sql``` path
* ensure the db details are updated correctly in ```local.yaml``` file in ```releases``` folder
* keep the ```local.yaml``` in the same folder as the executable
* go version 1.20 and above is needed to Build the project
    * clone the project to folder ```aspire-assignment```
    * run command ```go build -o aspire .``` for mac/linux
    * run command ```go build -o aspire.exe .``` for windows
    * run ```./aspire``` or ```aspire.exe```
    * the console should show a message ```starting router``` which means that the app has successfully started
* accesible psql database

---

