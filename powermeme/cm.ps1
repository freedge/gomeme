function Get-CMServer {
    param (
        [Parameter(Mandatory=$true)][PSCredential] $Credentials,
        [Parameter(Mandatory=$true)][string] $EndPoint,
        [Switch] $Insecure
    )

    $extraParams = @{
    }
    if ($Insecure) {
        $extraParams.add("-SkipCertificateCheck", $true)
    }

    $loginData = @{ 
        username = $credentials.UserName; 
        password = $credentials.GetNetworkCredential().Password}

    $res = Invoke-RestMethod -Method Post -Uri $endPoint/session/login  -Body (ConvertTo-Json $loginData) -ContentType "application/json" @extraParams
    
    return @{ Token=$res.token ;  EndPoint = $EndPoint ; ExtraParams = $extraParams}
}

function Get-CMJobStatus {
    param (
        [Parameter(Mandatory=$true)] $EmServer,
        [int] $Limit,
        [string] $Ctm,
        [string] $Name,
        [string] $Folder,
        [string] $Application,
        [string] $Id,
        [string] $HostName,
        [string] $Status
    )
    $headers = @{Authorization="Bearer "+ $EmServer.Token}
    $url = $EmServer.EndPoint + "/run/jobs/status"
    $extraparams = $EmServer.ExtraParams
    $query = @{}
    if ($PSBoundParameters.ContainsKey('Limit')) {
        $query.add("limit", $Limit) 
    }
    if ($PSBoundParameters.ContainsKey('Ctm')) {
        $query.add("ctm", $Ctm) 
    }
    if ($PSBoundParameters.ContainsKey('Name')) {
        $query.add("jobname", $Name) 
    }
    if ($PSBoundParameters.ContainsKey('Folder')) {
        $query.add("folder", $Folder) 
    }
    if ($PSBoundParameters.ContainsKey('Application')) {
        $query.add("application", $Application) 
    }
    if ($PSBoundParameters.ContainsKey('Id')) {
        $query.add("jobid", $Id) 
    }
    if ($PSBoundParameters.ContainsKey('HostName')) {
        $query.add("host", $HostName) 
    }
    if ($PSBoundParameters.ContainsKey('Status')) {
        $query.add("status", $Status) 
    }
    return (Invoke-RestMethod -Method Get -Uri $url -Headers $headers @extraparams -Body $query).Statuses
}

function Get-CMResource {
    param (
        [Parameter(Mandatory=$true)] $EmServer,
        [string] $Name,
        [string] $Ctm
    )
    $headers = @{Authorization="Bearer "+ $EmServer.Token}
    $url = $EmServer.EndPoint + "/run/resources"
    $extraparams = $EmServer.ExtraParams
    $query = @{}
    if ($PSBoundParameters.ContainsKey('Name')) {
        $query.add("name", $Name) 
    }
    if ($PSBoundParameters.ContainsKey('Ctm')) {
        $query.add("ctm", $Ctm) 
    }
    return Invoke-RestMethod -Method Get -Uri $url -Headers $headers @extraparams -Body $query
}


function Get-CMJobOutput {
    param (
        [Parameter(Mandatory=$true)] $EmServer,
        [Parameter(Mandatory=$true, ValueFromPipelineByPropertyName=$true)][Alias('Id')] [string[]] $jobId,
        [switch] $Output
    )
    $headers = @{Authorization="Bearer "+ $EmServer.Token; "Annotation-Subject"="Dummy"; "Annotation-Description"="Dummy"}
    $extraparams = $EmServer.ExtraParams
    $Option = "log"
    if ( $Output) {
        $Option = "output"
    }
    ForEach ( $input in $jobId)  {
        $url = $EmServer.EndPoint + "/run/job/" + $input + "/" + $Option
        Invoke-WebRequest -Uri $url -Headers $headers @extraparams | Write-Host
    }
}

