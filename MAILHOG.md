# Using MailHog with Craft

Nitro comes with [MailHog](https://github.com/mailhog/MailHog) for convenient local email testing and development.

## Set Up MailHog

MailHog requires setup for each machine you’d like to use it on. Run this command to spin it up:

```sh
nitro install mailhog
```

## Configure Craft Email for MailHog

MailHog’s ready to be used once it’s running, but it doesn’t change any of your mail settings by default. You can tell Craft or any app to send mail using MailHog’s SMTP settings.

From the Craft control panel, visit Settings → Email and enter the following:

- Transport Type: `SMTP`
- Host Name: `127.0.0.1`
- Port: `1025`
- Use Authentication: Unchecked (default)
- Encryption Method: `None` (default)
- Timeout: `10` (default)

## Access MailHog GUI

Once it’s running, the MailHog front end will be available at `http://[machine-ip]:8025`. After Craft is configured to use MailHog you should see the test email appear in its inbox.
