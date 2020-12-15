const express = require("express");
const {
  PermissionMiddlewareCreator,
  RecordSerializer,
} = require("forest-express-sequelize");
const models = require("../models");

const account = models.sequelize.account.models;
const webhookDb = models.sequelize.webhook.models;
const sequelize = models.sequelize;

const router = express.Router();
const permissionMiddlewareCreator = new PermissionMiddlewareCreator("account");

// This file contains the logic of every route in Forest Admin for the collection account:
// - Native routes are already generated but can be extended/overriden - Learn how to extend a route here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/extend-a-route
// - Smart action routes will need to be added as you create new Smart Actions - Learn how to create a Smart Action here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/actions/create-and-manage-smart-actions

// Create a Account
router.post(
  "/account",
  permissionMiddlewareCreator.create(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#create-a-record
    next();
  }
);

// Update a Account
router.put(
  "/account/:recordId",
  permissionMiddlewareCreator.update(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#update-a-record
    next();
  }
);

// Delete a Account
router.delete(
  "/account/:recordId",
  permissionMiddlewareCreator.delete(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#delete-a-record
    next();
  }
);

// Get a list of Accounts
router.get(
  "/account",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#get-a-list-of-records
    next();
  }
);

// Get a number of Accounts
router.get(
  "/account/count",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#get-a-number-of-records
    next();
  }
);

// Get a Account
router.get(
  "/account/:recordId",
  permissionMiddlewareCreator.details(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#get-a-record
    next();
  }
);

// Export a list of Accounts
router.get(
  "/account.csv",
  permissionMiddlewareCreator.export(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#export-a-list-of-records
    next();
  }
);

// Delete a list of Accounts
router.delete(
  "/account",
  permissionMiddlewareCreator.delete(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#delete-a-list-of-records
    next();
  }
);

// This needs to be here so stop errors on Account > Webhooks > Create
router.post(
  "/account/:recordId/relationships/webhooks",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    const accountId = request.params.recordId;
    const recordSerializer = new RecordSerializer(webhookDb.webhook);

    webhookDb.webhook
      .findAll({ where: { account_id: accountId } })
      .then((records) =>
        recordSerializer.serialize(records, { count: records.length })
      )
      .then((recordsSerialize) => response.send(recordsSerialize))
      .catch(next);
  }
);

router.get(
  "/account/:recordId/relationships/webhooks",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    const accountId = request.params.recordId;
    const recordSerializer = new RecordSerializer(webhookDb.webhook);

    webhookDb.webhook
      .findAll({ where: { account_id: accountId } })
      .then((records) =>
        recordSerializer.serialize(records, { count: records.length })
      )
      .then((recordsSerialize) => response.send(recordsSerialize))
      .catch(next);
  }
);

module.exports = router;
