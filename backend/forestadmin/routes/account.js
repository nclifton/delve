const express = require("express");
const {
  PermissionMiddlewareCreator,
  RecordSerializer,
} = require("forest-express-sequelize");
const { default: Axios } = require("axios");
const models = require("../models");

const account = models.sequelize.account.models;
const webhookDb = models.sequelize.webhook.models;
const sequelize = models.sequelize;

const router = express.Router();
const permissionMiddlewareCreator = new PermissionMiddlewareCreator("account");

var moment = require("moment");

var stringify = require("csv-stringify");

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
var CSV_OPTIONS = {
  formatters: {
    date: function date(value) {
      return moment(value).format();
    },
  },
};

const generateCSVResponse = (request, response, records) => {
  var filename = "".concat(request.params.filename, ".csv");
  response.setHeader("Content-Type", "text/csv; charset=utf-8");
  response.setHeader(
    "Content-disposition",
    "attachment; filename=".concat(filename)
  );
  response.setHeader("Last-Modified", moment());
  response.setHeader("X-Accel-Buffering", "no");
  response.setHeader("Cache-Control", "no-cache");
  var CSVAttributes = records.length > 0 ? Object.keys(records[0]) : [];
  var CSVLines = [];
  CSVLines.push(CSVAttributes);
  records.forEach((record) => {
    var CSVLine = [];
    CSVAttributes.forEach((attribute) => {
      var value;
      value = record[attribute];
      CSVLine.push(value);
    });
    CSVLines.push(CSVLine);
  });
  stringify(CSVLines, CSV_OPTIONS, function (error, csv) {
    response.write(csv);
    response.end();
  });
};

router.post(
  "/actions/export-global-usage-report",
  permissionMiddlewareCreator.list(),
  async (request, response, next) => {
    let report = {};
    try {
      report = await Axios({
        url: `${process.env.FOREST_ADMIN_API}/report/usage`,
        method: "get",
        headers: {
          "content-type": "application/json",
        },
      });
    } catch (err) {
      next(err);
    }

    generateCSVResponse(request, response, report.data);
  }
);

router.post(
  "/actions/export-account-usage-report",
  permissionMiddlewareCreator.list(),
  async (request, response, next) => {
    let report = {};
    try {
      requests = [];
      for (const id of request.body.data.attributes.ids) {
        requests.push(
          Axios({
            url: `${process.env.FOREST_ADMIN_API}/report/usage/${id}`,
            method: "get",
            headers: {
              "content-type": "application/json",
            },
          })
        );
      }
      values = await Promise.all(requests);
      report = values.map((x) => {
        return x.data.length > 0 ? x.data[0] : {};
      });
    } catch (err) {
      next(err);
    }

    generateCSVResponse(request, response, report);
  }
);

module.exports = router;
