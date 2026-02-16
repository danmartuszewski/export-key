package shell

const FishInit = `# export-key shell integration
function ek
  set -l output (export-key select $argv)
  set -l exit_code $status
  if test $exit_code -eq 0; and test -n "$output"
    eval $output
  end
  return $exit_code
end
`
