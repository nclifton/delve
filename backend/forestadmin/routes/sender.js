const express = require("express");
const {
  PermissionMiddlewareCreator,
  RecordCreator,
  RecordUpdater,
  RecordSerializer,
} = require("forest-express-sequelize");
const models = require("../models");
const parseDataUri = require('parse-data-uri');
const csv = require('csv');
const P = require('bluebird');

const sender = models.sequelize.sender.models;
const dbAccount = models.sequelize.account.models;

const router = express.Router();
const permissionMiddlewareCreator = new PermissionMiddlewareCreator("sender");

// This file contains the logic of every route in Forest Admin for the collection sender:
// - Native routes are already generated but can be extended/overriden - Learn how to extend a route here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/extend-a-route
// - Smart action routes will need to be added as you create new Smart Actions - Learn how to create a Smart Action here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/actions/create-and-manage-smart-actions

// Create a Sender
router.post(
  "/sender",
  permissionMiddlewareCreator.create(),
  (request, response, next) => {
    const recordCreator = new RecordCreator(sender.sender);
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
  "/sender/:recordId",
  permissionMiddlewareCreator.update(),
  (request, response, next) => {
    const recordUpdater = new RecordUpdater(sender.sender);
    recordUpdater
      .deserialize(request.body)
      .then((recordToUpdate) => {
        if (request.body.data.relationships.account_id.data) {
          recordToUpdate.ref_account_id =
            request.body.data.relationships.account_id.data.id;
        } else {
          recordToUpdate.ref_account_id = null
        }
        return recordUpdater.update(recordToUpdate, request.params.recordId);
      })
      .then((record) => recordUpdater.serialize(record))
      .then((recordSerialized) => response.send(recordSerialized))
      .catch(next);
  }
);

// Delete a Account
router.delete(
  "/sender/:recordId",
  permissionMiddlewareCreator.delete(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#delete-a-record
    next();
  }
);

// Get a list of Accounts
router.get(
  "/sender",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#get-a-list-of-records
    next();
  }
);

// Get a number of Accounts
router.get(
  "/sender/count",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#get-a-number-of-records
    next();
  }
);

// Get a Account
router.get(
  "/sender/:recordId",
  permissionMiddlewareCreator.details(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#get-a-record
    next();
  }
);

// Export a list of Accounts
router.get(
  "/sender.csv",
  permissionMiddlewareCreator.export(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#export-a-list-of-records
    next();
  }
);

// Delete a list of Accounts
router.delete(
  "/sender",
  permissionMiddlewareCreator.delete(),
  (request, response, next) => {
    // Learn what this route does here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/routes/default-routes#delete-a-list-of-records
    next();
  }
);

// This needs to be here to make updating account_id on sender work..
router.put(
  "/sender/:recordId/relationships/account_id",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    const senderId = request.params.recordId;
    const recordUpdater = new RecordUpdater(sender.sender);

    sender.sender
      .findByPk(senderId)
      .then((senderRecord) => {
        recordToUpdate = {
          id: senderId,
          ref_account_id: senderRecord.ref_account_id,
          account_id: senderRecord.ref_account_id,
        };
        return recordUpdater.update(recordToUpdate, senderId);
      })
      .then((record) => recordUpdater.serialize(record))
      .then((recordSerialized) => response.send(recordSerialized))
      .catch(next);
  }
);

router.post(
  '/sender/import',
  permissionMiddlewareCreator.create(),
  (req, res, next) => {
    let parsed = parseDataUri(req.body.data.attributes.values['CSV file']);

    //TODO need to deal with duplicates
    //TODO field validation
    //TODO errors output
    //TODO account name?

    csv.parse(parsed.data, { columns: true }, function (err, rows) {
      if (err) {
        res.status(400).send({
          error: `Cannot import data: ${err.message}`
        });
      } else {
        P.each(rows, (row => {

          const recordCreator = new RecordCreator(sender.sender)
          recordCreator.create({
            address: row.address,
            country: row.country,
            comment: row.comment,
            channels: JSON.parse(row.channels),
            mmsProviderKey: row.mms_provider_key == "" ? null : row.mms_provider_key,
            account_id: row.account_id == "" ? null : row.account_id,
            ref_account_id: row.account_id == "" ? null : row.account_id,
          });

        })).then(() => {
          res.send({ success: 'Data successfully imported!' });
        })
          .catch(next);
      }
    });
  });

module.exports = router;
