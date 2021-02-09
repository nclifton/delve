const { collection } = require("forest-express-sequelize");
const models = require("../models");

const dbAccount = models.sequelize.account.models;

// This file allows you to add to your Forest UI:
// - Smart actions: https://docs.forestadmin.com/documentation/reference-guide/actions/create-and-manage-smart-actions
// - Smart fields: https://docs.forestadmin.com/documentation/reference-guide/fields/create-and-manage-smart-fields
// - Smart relationships: https://docs.forestadmin.com/documentation/reference-guide/relationships/create-a-smart-relationship
// - Smart segments: https://docs.forestadmin.com/documentation/reference-guide/segments/smart-segments
collection("sender", {
  actions: [{
    name: 'Import senders',
    endpoint: '/forest/sender/import',
    type: 'global',
    fields: [{
      field: 'CSV file',
      description: 'A comma separated values file stores tabular data (numbers and text) in plain text',
      type: 'File',
      isRequired: true
    }]
  }],
  fields: [
    {
      field: "account_id",
      type: "String",
      reference: "account.id",
      get: function (sender) {
        return dbAccount.account.findByPk(sender.ref_account_id);
      },
    },
  ],
  segments: [],
});
