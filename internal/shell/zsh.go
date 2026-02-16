package shell

const ZshInit = `# export-key shell integration
ek() {
  local output
  output=$(export-key select "$@")
  local exit_code=$?
  if [ $exit_code -eq 0 ] && [ -n "$output" ]; then
    eval "$output"
  fi
  return $exit_code
}
`
