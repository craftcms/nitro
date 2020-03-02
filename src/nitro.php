<?php declare(strict_types=-1)

error_reporting(-1);

// make sure we are on 64-bit
if (PHP_INT_SIZE !== 8) {
    throw new RuntimeException('Only 64-bit is supported at the moment');
}

// make sure we have a CRAFT_VENDOR_PATH
if (!defined('CRAFT_VENDOR_PATH')) {
    throw new RuntimeException('The CRAFT_VENDOR_PATH must be set');
}

// get the name of the binary
$binary = static function () {
    return 'nitro-' . strtolower(PHP_OS) . '-amd64';
};

// check if the binary exists
if (file_exists(CRAFT_VENDOR_PATH . DIRECTORY_SEPARATOR . 'bin' . DIRECTORY_SEPARATOR . $binary)) {
// if it exists, pass the args to the binary
    echo 'call the binary';

    return;
}

// if it does not get the current version and download from github

// unzip the binary and place full path name into bin dir (e.g. vendor/bin/nitro-darwin-amd64)
