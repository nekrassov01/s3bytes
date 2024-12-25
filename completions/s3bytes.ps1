Register-ArgumentCompleter -Native -CommandName "s3bytes" -ScriptBlock {
    param($commandName, $wordToComplete, $cursorPosition)
    (Invoke-Expression "$wordToComplete --generate-bash-completion").ForEach{
        [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}
