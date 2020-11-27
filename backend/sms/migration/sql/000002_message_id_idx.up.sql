BEGIN;
CREATE INDEX sms_message_id ON sms(message_id);
CREATE INDEX sms_related_to_mo ON sms(account_id, sender, recipient);
COMMIT;
