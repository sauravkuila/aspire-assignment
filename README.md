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
* API version management put in place for ease of management as product grows

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

## API Endpoints
The postman collection in ```releases/aspire-assignment.postman_collection.json``` will ensure all APIs are documented with relevant tests to sync tokens in collection variables

* `GET`    /health                   --> health check api. can be used for k8 pod health or circuit breaker
* `POST`   /cred/signup              --> signup api. works without any auth
* `POST`   /cred/login               --> login api. works without any auth
* `POST`   /v1/loan                  --> apply loan api. only authenticated customer can reach this
* `PUT`    /v1/loan                  --> modify loan api. only authenticated customer can reach this
* `DELETE` /v1/loan                  --> cancel loan api. only authenticated customer can reach this
* `GET`    /v1/loan/status           --> get loan status. only authenticated customer can reach this
* `GET`    /v1/loan/installments     --> get loan installments and their status. only authenticated customer can reach this
* `POST`   /v1/loan/repay            --> customer scheduled payment api. only authenticated customer can reach this
* `GET`    /v1/admin/applications    --> lists pending loans. only authenticated admin can reach this
* `POST`   /v1/admin/update          --> approve/reject pending loans. only authenticated admin can reach this

### Usage
* Download the relevant executable from `releases/macos` or `releases/windows` folder and run
* Download the `local.yaml` and edit the database connection settings
```
databases:
  postgres:
    host: 127.0.0.1     #db connection ip
    port: 5432          #db connection port
    user: postgres      #db username
    password: postgres  #db password
    db: aspire          #db name
    sslmode: disable
    connect_timeout: 10
```
* Run the executable ```./aspire```(mac) or ```aspire.exe```(windows)
    * the console should show a message ```starting router``` which means that the app has successfully started
    * ensure to download the `local.yaml` and keep it in the same folder as the executable
* Import the Postman collection from ```releases/aspire-assignment.postman_collection.json```
* Signup using `/cred/signup` and create a username and password as a `CUTOMER` or `ADMIN`
* Login using `/cred/login` and receive a auth token to be used for all loan APIs
* Apply for a loan using `/v1/loan`
* Check loan status using `/v1/loan/status`
* Login as an `ADMIN` and check if loan application is available for approve/reject using `/v1/admin/applications`
* As an `ADMIN`, approve the loan using `/v1/admin/update`
* Login as the initial user and check the loan status using `/v1/loan/status`
    * If the loan is approved, the loan state will show `APPROVED` and the installments will show as `PENDING` in `/v1/loan/installments`
    * If the loan is rejected, the loan state will show `REJECTED` and the installments will not show in `/v1/loan/installments`
* Pay a loan installment using `/v1/loan/repay`
    * Installment amount less than amount due will not be accepted
    * installment amount greater than amount due will be accepted and the upcoming payments will be recalculated. the same can be observed with `/v1/loan/installments` after each payment
    * payments mark the scheduled payment as `PAID`
    * The loan is marked as `PAID` when the ourstanding amount in `/v1/loan/installments` response becomes 0
    * If the loan is repayed before scheduled tenure, the remaining payments are marked `CANCELLED`

---

Happy Usage!
