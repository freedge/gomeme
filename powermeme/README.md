Based on https://github.com/controlm/automation-api-quickstart/blob/master/201-call-rest-api-using-powershell/AutomationAPIExample.ps1

Requires at least PowerShell 6

```

$server = Get-CMServer -Endpoint "https://workbench:8443/automation-api" -Insecure
Get-CMJobStatus $server -Application *A* | Out-GridView -PassThru | Get-CMJobOutput $server -Output
```