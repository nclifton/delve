const express = require("express");
const {
  PermissionMiddlewareCreator,
  RecordSerializer,
  RecordUpdater,
} = require("forest-express-sequelize");
const { default: Axios } = require("axios");
const models = require("../models");

const account = models.sequelize.account.models;
const webhookDb = models.sequelize.webhook.models;
const senderDb = models.sequelize.sender.models;
const smsDb = models.sequelize.sms.models;
const mmsDb = models.sequelize.mms.models;
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
    const limit = parseInt(request.query.page.size, 10) || 20;
    const offset = (parseInt(request.query.page.number, 10) - 1) * limit;
    const recordSerializer = new RecordSerializer(webhookDb.webhook);

    let countQuery = `
      SELECT count(*)
      FROM webhook
      WHERE account_id = '${accountId}'
    `;

    let findQuery = `
      SELECT webhook.*
      FROM webhook
      WHERE account_id = '${accountId}'
    `;

    if (request.query.search) {
      filter = `AND (
        name ILIKE '%${request.query.search}%'
        OR event ILIKE '%${request.query.search}%'
        OR url ILIKE '%${request.query.search}%'
      )`;
      countQuery += filter;
      findQuery += filter;
    }

    findQuery += `
      LIMIT ${limit}
      OFFSET ${offset}
    `;

    const countqry = sequelize.webhook.query(countQuery, {
      type: sequelize.webhook.QueryTypes.SELECT,
    });
    const findqry = sequelize.webhook.query(findQuery, {
      type: sequelize.webhook.QueryTypes.SELECT,
    });

    Promise.all([countqry, findqry])
      .then(([count, records]) =>
        recordSerializer.serialize(records, { count: count[0].count })
      )
      .then((recordsSerialize) => response.send(recordsSerialize))
      .catch(next);
  }
);

router.post(
  "/account/:recordId/relationships/senders",
  permissionMiddlewareCreator.create(),
  (request, response, next) => {
    const accountId = request.params.recordId;
    const recordSerializer = new RecordSerializer(senderDb.sender);

    senderDb.sender
      .findAll({ where: { account_id: accountId } })
      .then((records) =>
        recordSerializer.serialize(records, { count: records.length })
      )
      .then((recordsSerialize) => response.send(recordsSerialize))
      .catch(next);
  }
);

router.delete(
  "/account/:recordId/relationships/senders",
  permissionMiddlewareCreator.delete(),
  (request, response, next) => {
    const accountId = request.params.recordId;
    const recordUpdater = new RecordUpdater(senderDb.sender);

    request.body.data.forEach((selected)=>{
      senderDb.sender
      .findAll({ where: { id: selected.id, account_id: accountId } })
      .then(() => {
        recordToUpdate = {
          id: selected.id,
          ref_account_id: null,
          account_id: null,
        };
        return recordUpdater.update(recordToUpdate, selected.id);
      })
      .catch(next);
    })

  }
);

router.get(
  "/account/:recordId/relationships/senders",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    const accountId = request.params.recordId;
    const limit = parseInt(request.query.page.size, 10) || 20;
    const offset = (parseInt(request.query.page.number, 10) - 1) * limit;
    const recordSerializer = new RecordSerializer(senderDb.sender);

    let countQuery = `
      SELECT count(*)
      FROM sender
      WHERE account_id = '${accountId}'
    `;

    let findQuery = `
      SELECT sender.*
      FROM sender
      WHERE account_id = '${accountId}'
    `;

    if (request.query.search) {
      filter = `AND (
        address ILIKE '%${request.query.search}%'
        OR country ILIKE '%${request.query.search}%'
        OR channels ILIKE '%${request.query.search}%'
        OR mms_provider_key ILIKE '%${request.query.search}%'
      )`;
      countQuery += filter;
      findQuery += filter;
    }

    findQuery += `
      LIMIT ${limit}
      OFFSET ${offset}
    `;

    const countqry = sequelize.sender.query(countQuery, {
      type: sequelize.sender.QueryTypes.SELECT,
    });
    const findqry = sequelize.sender.query(findQuery, {
      type: sequelize.sender.QueryTypes.SELECT,
    });

    Promise.all([countqry, findqry])
      .then(([count, records]) =>
        recordSerializer.serialize(records, { count: count[0].count })
      )
      .then((recordsSerialize) => response.send(recordsSerialize))
      .catch(next);
  }
);

router.get(
  "/account/:recordId/relationships/sms",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    const accountId = request.params.recordId;
    const limit = parseInt(request.query.page.size, 10) || 20;
    const offset = (parseInt(request.query.page.number, 10) - 1) * limit;
    const recordSerializer = new RecordSerializer(smsDb.sms);

    let countQuery = `
      SELECT count(*)
      FROM sms
      WHERE account_id = '${accountId}'
    `;

    let findQuery = `
      SELECT sms.*
      FROM sms
      WHERE account_id = '${accountId}'
    `;

    if (request.query.search) {
      filter = `AND (
        ARRAY[message_ref, status, message_id, sender, recipient, country] && ARRAY[LOWER('${request.query.search}'), UPPER('${request.query.search}')]
        OR message ILIKE '%${request.query.search}%'
      )`;
      countQuery += filter;
      findQuery += filter;
    }

    findQuery += `
      LIMIT ${limit}
      OFFSET ${offset}
    `;

    const countqry = sequelize.sms.query(countQuery, {
      type: sequelize.sms.QueryTypes.SELECT,
    });
    const findqry = sequelize.sms.query(findQuery, {
      type: sequelize.sms.QueryTypes.SELECT,
    });

    Promise.all([countqry, findqry])
      .then(([count, records]) =>
        recordSerializer.serialize(records, { count: count[0].count })
      )
      .then((recordsSerialize) => response.send(recordsSerialize))
      .catch(next);
  }
);

router.get(
  "/account/:recordId/relationships/mms",
  permissionMiddlewareCreator.list(),
  (request, response, next) => {
    const accountId = request.params.recordId;
    console.log(request.query);
    const limit = parseInt(request.query.page.size, 10) || 20;
    const offset = (parseInt(request.query.page.number, 10) - 1) * limit;
    const recordSerializer = new RecordSerializer(mmsDb.mms);

    let countQuery = `
      SELECT count(*)
      FROM mms
      WHERE account_id = '${accountId}'
    `;

    let findQuery = `
      SELECT mms.*
      FROM mms
      WHERE account_id = '${accountId}'
    `;

    if (request.query.search) {
      filter = `AND (
        ARRAY[message_ref, status, message_id, sender, recipient, country] && ARRAY[LOWER('${request.query.search}'), UPPER('${request.query.search}')]
        OR subject ILIKE '%${request.query.search}%'
        OR message ILIKE '%${request.query.search}%'
      )`;
      countQuery += filter;
      findQuery += filter;
    }

    findQuery += `
      LIMIT ${limit}
      OFFSET ${offset}
    `;

    const countqry = sequelize.mms.query(countQuery, {
      type: sequelize.mms.QueryTypes.SELECT,
    });
    const findqry = sequelize.mms.query(findQuery, {
      type: sequelize.mms.QueryTypes.SELECT,
    });

    Promise.all([countqry, findqry])
      .then(([count, records]) =>
        recordSerializer.serialize(records, { count: count[0].count })
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
