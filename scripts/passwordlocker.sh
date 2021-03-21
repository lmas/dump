#!/bin/bash -e
#
# Extra depencies (besides standard tools like shasum, tee etc.):
# - gpg
# - xsel
# - pwgen

USERID="tester@tester.se"
PASSDIR="./pass"
CLIPBOARD_TIMEOUT="5"

################################################################################

# Decrypt and show the whole file (first line is the password).
function cmdShow() {
    pass="$(gpg --quiet --armor --decrypt "$PASSDIR/$fileid")"
    echo "$pass"
}

# Decrypt and paste the password (from the first line) into the clipboard.
# Clipboard is cleared automagically after some time.
function cmdClip() {
    # read grabs only the first line from the decrypted file
    read pass <<< "$(gpg --quiet --armor --decrypt "$PASSDIR/$fileid")"
    echo -n "$pass" | xsel --input --primary # TODO: not sure we should use the primary selection?
    (
        sleep "$CLIPBOARD_TIMEOUT"
        xsel --clear --primary
    ) & disown # prevent getting any SIGHUP and makes sure we clear the clipboard
}

# Read in a new password (and any extra info after the first line) from STDIN
# and save it, possibly overwriting an existing file.
function cmdSave() {
    echo "Enter new password to save (ctrl+d to stop):"
    pass="$(tee)"
    echo "$pass" | gpg --quiet --armor --recipient "$USERID" --output "$PASSDIR/$fileid" --encrypt
}

# Remove a password file.
function cmdDelete() {
    rm "$PASSDIR/$fileid"
}

# Generate a completely new password and save it, possibly overwriting an existing file.
function cmdGenerate() {
    pass="$(pwgen -s 16 1)"
    echo "$pass" | gpg --quiet --armor --recipient "$USERID" --output "$PASSDIR/$fileid" --encrypt
}

################################################################################

cmd="$1"
fileid=$(echo "$2" | sha256sum - | cut -d' ' -f1)

case "$cmd" in
    show)
        cmdShow ;;
    clip)
        cmdClip ;;
    save)
        cmdSave ;;
    del)
        cmdDelete ;;
    gen)
        cmdGenerate ;;
    help)
        cmdUsage ;;
    *)
        echo "Unknown command: $cmd"; exit 1 ;;
esac

