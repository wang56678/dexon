/^import (/a\
goerrors "errors"
s|errors\.New(|goerrors.New(|g
s|\([^[:alnum:]]\)Parse\([[:upper:]]\{0,1\}[[:lower:]]*\)\([( ]\)|\1parse\2\3|g
