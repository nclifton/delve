const { collection } = require("forest-express-sequelize");

// This file allows you to add to your Forest UI:
// - Smart actions: https://docs.forestadmin.com/documentation/reference-guide/actions/create-and-manage-smart-actions
// - Smart fields: https://docs.forestadmin.com/documentation/reference-guide/fields/create-and-manage-smart-fields
// - Smart relationships: https://docs.forestadmin.com/documentation/reference-guide/relationships/create-a-smart-relationship
// - Smart segments: https://docs.forestadmin.com/documentation/reference-guide/segments/smart-segments
collection("account", {
  actions: [
    {
      name: "Export global usage report",
      type: "global",
      download: true,
    },
    {
      name: "Export account usage report",
      type: "bulk",
      download: true,
    },
  ],
  fields: [
    {
      field: "webhooks",
      type: ["String"],
      reference: "webhook.account_id",
    },
    {
      field: "mms",
      type: ["String"],
      reference: "mms.account_id",
    },
    {
      field: "sms",
      type: ["String"],
      reference: "sms.account_id",
    },
  ],
  segments: [],
});
