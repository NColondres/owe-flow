Google doc about auth https://developers.google.com/workspace/guides/configure-oauth-consent
svc account https://developers.google.com/workspace/guides/create-credentials#service-account

google recommends auth using https://cloud.google.com/iam/docs/workload-identity-federation


Created a SVC account in google, generated a key and credentials for it.
You will need to share the google sheets file with the svc account. Use the email address provided with the svc account.

Code checks for env vars. It looks for two.
1. is the json file for the svc account
2. the id of the google sheets file.

You can get the google sheets id in the browser. 

From the URL: The Spreadsheet ID is the string of characters between "/d/" and "/edit" in the URL of your spreadsheet. For example, in the URL "https://docs.google.com/spreadsheets/d/1QPvcIcmNU1QbZYRF__rrjIC4C1F0Ir3KI-YtIRCCWws/edit#gid=0", the Spreadsheet ID is "1QPvcIcmNU1QbZYRF__rrjIC4C1F0Ir3KI-YtIRCCWws"

you'll need to run 
export SVC_PATH='path/to/jsonfile'
export GSHEET="GoogleSheetID"