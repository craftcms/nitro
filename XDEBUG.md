# Using Xdebug with Nitro and PhpStorm

Once you’ve created a machine, you can run `nitro xdebug on` and `nitro xdebug configure` to enable [Xdebug](https://xdebug.org/) and set it up for automatic connections from your host machine.

You can use `nitro xdebug off` to disable Xdebug without having to restart the machine.

## Configuring PhpStorm

First, you’ll need to configure PhpStorm to listen for requests from the browser or console.

1. Create a new server in PhpStorm using your machine’s domain name. (“Preferences” → “Languages & Frameworks” → “PHP” → “Servers”.)  
![PhpStorm Server Settings](resources/phpstorm_server.png)

2. Enable “Use path mappings” and set your existing project root to the absolute path on the server. The absolute path will look like `/home/ubuntu/sites/my-site`, where `my-site` reflects your actual project’s folder name in the Nitro machine. (Use `nitro context` if you need to check the path, and keep in mind this is the project root and not necessarily the web root.)

3. Choose “Run” → “Edit Configurations...” and create a new “PHP Remote Debug” configuration, selecting the server you just created. Check “Filter debug connection by IDE key” and enter `PHPSTORM`.  
![PhpStorm Remote Debug Settings](resources/phpstorm_remote_debug.png)

4. Choose “Run” → “Start Listening for PHP Debug Connections”.  
![PhpStorm Remote Debug Settings](resources/start_listening.png)

## Debugging Web Requests

1. Install the Xdebug helper in your favorite browser.

- [Chrome](https://chrome.google.com/extensions/detail/eadndfjplgieldjbigjakmdgkmoaaaoc)
- [Firefox](https://addons.mozilla.org/en-US/firefox/addon/xdebug-helper-for-firefox/)
- [Internet Explorer](https://www.jetbrains.com/phpstorm/marklets/)
- [Safari](https://github.com/benmatselby/xdebug-toggler)
- [Opera](https://addons.opera.com/addons/extensions/details/xdebug-launcher/)

2. In the browser helper’s options, choose “PhpStorm” and save.  
![Xdebug Browser Helper Chrome](resources/xdebug_chrome_settings.png)

3. Choose “Debug” on your browser’s Xdebug helper.  
![PhpStorm Remote Debug Settings](resources/xdebug_chrome.png)

4. Load the site in your browser and whatever breakpoints you’ve set will be hit.

## Debugging PHP console requests

SSH into your Nitro machine using `nitro ssh`, then run your PHP script from the console and any breakpoints you’ve set will be hit.
