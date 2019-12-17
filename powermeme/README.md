Based on https://github.com/controlm/automation-api-quickstart/blob/master/201-call-rest-api-using-powershell/AutomationAPIExample.ps1

Requires at least PowerShell 6

```
$creds = Get-Credential
$server = Get-CMServer -Credentials $creds -Endpoint "https://workbench:8443/automation-api" -Insecure
(Get-CMJobStatus $server).Statuses  | Out-GridView
```