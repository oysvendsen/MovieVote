main() {
  go test endpointtests/*.go -v
}

if [[ $(basename $(pwd)) == "MovieVote" ]]; then
  main
else
  echo "Can only execute endpoint tests from root-directory"
fi
