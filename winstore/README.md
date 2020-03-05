The story is that we have now our pipeline that generates a msix file and upload
it in the released files of github, which is great.

To enable automatic updates, we need as well, an appinstaller file. When downloading that 
appinstaller, the installer on Windows (for example, Add-AppxPackage -AppInstallerFile)
will read the Uri from the appinstaller file, then issue a HEAD request to the server.

Sadly, Github releases are served using signed URLs and HEAD request gives a 403.

So we need to serve files from somewhere else. For that we can use for example, the deploymentconfig from here. It's a bit overkill (we could simply upload to S3 for example)
but it's more fun.