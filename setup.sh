#!/bin/sh
# shellcheck disable=SC3043

#
#  Setup environment
#

export ORAS_BEGIN_ENV="$HOME/.config/oras"
export ORAS_NAMESPACE=${ORAS_NAMESPACE:=$1}

ORAS_ROOT=$ORAS_BEGIN_ENV/$ORAS_NAMESPACE

ORAS_LOGIN_DIR="$ORAS_ROOT/login"
ORAS_ACCESS_DIR="$ORAS_ROOT/access"
ORAS_STORE_DIR="$ORAS_ROOT/store"

mkdir -p "$ORAS_BEGIN_ENV"
mkdir -p "$ORAS_LOGIN_DIR"
mkdir -p "$ORAS_ACCESS_DIR"
mkdir -p "$ORAS_STORE_DIR"

#
#  Setup dependencies
#

ORAS_SOURCE_DIR=${ORAS_SOURCE_DIR:-"$HOME/azorasrc/src/rc"}
ORAS_RC_PATH="$ORAS_SOURCE_DIR/orasrc"
ORAS_ACCESS_PATH="$ORAS_SOURCE_DIR/access.alias"
ORAS_ENCPASS_PATH="$ORAS_SOURCE_DIR/encpass.sh"

if cp "$ORAS_RC_PATH" "$ORAS_ROOT/orasrc"; then
    ORAS_RC_PATH="$ORAS_BEGIN_ENV/$ORAS_NAMESPACE/orasrc"
fi

if cp "$ORAS_ACCESS_PATH" "$ORAS_ROOT/access.alias"; then 
    ORAS_ACCESS_PATH="$ORAS_ROOT/access.alias"
fi

if cp "$ORAS_ENCPASS_PATH" "$ORAS_ROOT/encpass.sh"; then
    ORAS_ENCPASS_PATH="$ORAS_ROOT/encpass.sh"
fi

#
# Configure components for namespace
#

ORAS_LOGIN_RC="$ORAS_LOGIN_DIR/loginrc"
ORAS_ACCESS_RC="$ORAS_ACCESS_DIR/accessrc"
ORAS_STORE_RC="$ORAS_STORE_DIR/storerc"

: > "$ORAS_LOGIN_RC"
: > "$ORAS_ACCESS_RC"
: > "$ORAS_STORE_RC"

cat <<-EOF > "$ORAS_LOGIN_RC"
#!/bin/sh

export ORAS_SOURCE_DIR="$ORAS_SOURCE_DIR"
export ORAS_NAMESPACE="$ORAS_NAMESPACE"
export ORAS_BEGIN_ENV="$ORAS_BEGIN_ENV"

_login() {
    ORAS_ACCESS_PATH="$ORAS_ACCESS_PATH" "$ORAS_RC_PATH" 'login' \$*
}

alias login=_login

if [ "\$#" -gt 0 ]; then
    login "\$@" 
fi

EOF

cat <<-EOF > "$ORAS_ACCESS_RC"
#!/bin/sh

export ORAS_SOURCE_DIR="$ORAS_SOURCE_DIR"
export ORAS_NAMESPACE="$ORAS_NAMESPACE"
export ORAS_BEGIN_ENV="$ORAS_BEGIN_ENV"

_access() {
    ORAS_ACCESS_PATH="$ORAS_ACCESS_PATH" $ORAS_RC_PATH 'access' \$*
}

alias access=_access

if [ "\$#" -gt 0 ]; then
    access "\$@" 
fi

EOF

cat <<-EOF > "$ORAS_STORE_RC"
#!/bin/sh

export ORAS_SOURCE_DIR="$ORAS_SOURCE_DIR"
export ORAS_NAMESPACE="$ORAS_NAMESPACE"
export ORAS_BEGIN_ENV="$ORAS_BEGIN_ENV"

_store() {
    ORAS_ACCESS_PATH="$ORAS_ACCESS_PATH" $ORAS_RC_PATH 'store' \$*
}
alias store=_store

if [ "\$#" -gt 0 ]; then
    store "\$@" 
fi

EOF

chmod +x "$ORAS_LOGIN_RC"
chmod +x "$ORAS_ACCESS_RC"
chmod +x "$ORAS_STORE_RC"

# shellcheck disable=SC1090
. "$ORAS_LOGIN_RC" > /dev/null
# shellcheck disable=SC1090
. "$ORAS_ACCESS_RC" > /dev/null
# shellcheck disable=SC1090
. "$ORAS_STORE_RC" > /dev/null
