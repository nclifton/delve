const express = require("express");
const {
  PermissionMiddlewareCreator,
  RecordCreator,
  RecordUpdater,
  RecordSerializer,
} = require("forest-express-sequelize");
const models = require("../models");

const webhook = models.sequelize.webhook.models;
const dbAccount = models.sequelize.account.models;

const router = express.Router();
const permissionMiddlewareCreator = new PermissionMiddlewareCreator("webhook");

// This file contains the logic of every route in Forest Admin for the collection webhook:
// - Native routes are already generated but can be extended/overriden - Learn how to extend a route here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/extend-a-route
// - Smart action routes will need to be added as you create new Smart Actions - Learn how to create a Smart Action here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/actions/create-and-manage-smart-actions

// Create a Account
router.post(
  "/webhook",
  permissionMiddlewareCreator.create(),
  (request, response, next) => {
    const recordCreator = new RecordCreator(webhook.webhook);
    recordCreator
      .deserialize(request.body)
      .then((recordToCreate) => {
        recordToCreate.ref_account_id = recordToCreate.account_id;
        return recordCreator.create(recordToCreate);
      })
      .then((record) => recordCreator.serialize(record))
      .then((recordSerialized) => response.send(recordSerialized))
      .catch(next);
  }
);

// Update a Account
router.put(
  "/webhook/:recordId",
  permissionMiddlewareCreator.update(),
  (request, response, next) => {
    const recordUpdater = new RecordUpdater(webhook.webhook);
    recordUpdater
      .deserialize(request.body)
      .then((recordToUpdate) => {
        recordToUpdate.ref_account_id =
          request.body.data.relationships.account_id.data.id;
        return recordUpdater.update(recordToUpdate, request.params.recordId);
      })
      .then((record) => recordUpdater.serialize(record))
      .then((recordSerialized) => response.send(recordSerialized))
      .catch(next);
  }
);

// Delete a Account
router.delete(
  "/webhook/:recordId",
  permissionMiddlewareCreator.delete(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#delete-a-record
    next();
  }
);

// Get a list of Accounts
router.get(
  "/webhook",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#get-a-list-of-records
    next();
  }
);

// Get a number of Accounts
router.get(
  "/webhook/count",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#get-a-number-of-records
    next();
  }
);

// Get a Account
router.get(
  "/webhook/:recordId",
  permissionMiddlewareCreator.details(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#get-a-record
    next();
  }
);

// Export a list of Accounts
router.get(
  "/webhook.csv",
  permissionMiddlewareCreator.export(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#export-a-list-of-records
    next();
  }
);

// Delete a list of Accounts
router.delete(
  "/webhook",
  permissionMiddlewareCreator.delete(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#delete-a-list-of-records
    next();
  }
);

// This needs to be here to make updating account_id on webhook work..
router.put(
  "/webhook/:recordId/relationships/account_id",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    const webhookId = request.params.recordId;
    const recordUpdater = new RecordUpdater(webhook.webhook);

    webhook.webhook
      .findByPk(webhookId)
      .then((webhookRecord) => {
        recordToUpdate = {
          id: webhookId,
          ref_account_id: webhookRecord.ref_account_id,
          account_id: webhookRecord.ref_account_id,
        };
        return recordUpdater.update(recordToUpdate, webhookId);
      })
      .then((record) => recordUpdater.serialize(record))
      .then((recordSerialized) => response.send(recordSerialized))
      .catch(next);
  }
);

module.exports = router;
