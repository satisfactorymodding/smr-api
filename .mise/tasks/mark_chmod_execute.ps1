#MISE description="For developers on Windows, mark all .sh file-based Mise tasks as executable so Linux devs and the CI can use them"

$shFiles = Get-ChildItem -Path "./.mise/tasks" -Filter "*.sh" -File

foreach ($file in $shFiles) {
    Write-Host "Setting Linux execute permission on $($file.Name)"
    # https://git-scm.com/docs/git-update-index#Documentation/git-update-index.txt---chmod-x
    git update-index --add --chmod=+x $file.FullName
}
