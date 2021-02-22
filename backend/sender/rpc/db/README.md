# Sender DB


## Mock Generation

Generated mock for the DB for use in SQL unit (integration) tests.

mock DB (`db.MockDB`) is generated using [ vektra/mockery ](https://github.com/vektra/mockery/blob/master/pkg/generator.go).

command line:

```
cd backend/sender/rpc/db
mockery --name DB --inpackage
```