# Using Blackfire with Nitro

Once you have a machine created, you can configure and run [Blackfire](https://blackfire.io/).

View the [full installation docs](https://blackfire.io/docs/up-and-running/installation?action=install&mode=full&location=local&os=debian&language=php).

> If XDebug is installed then it may alter performance measures significantly. You can disable it as follows:
>```
>nitro xdebug off
>```

## With Browser Extensions
 
1. Log in or sign up for an account at [https://blackfire.io/](https://blackfire.io/).

2. SSH into Nitro: 
```
nitro ssh
```

3. Configure the local agent with the server credentials ([https://blackfire.io/my/settings/credentials](https://blackfire.io/my/settings/credentials)): 
```
sudo blackfire-agent --register
```

4. Restart the agent service: 
```
sudo /etc/init.d/blackfire-agent restart
```

5. Install the [Firefox extension](https://addons.mozilla.org/en-GB/firefox/addon/blackfire/) or the [Chrome extension](https://chrome.google.com/webstore/detail/blackfire-profiler/miefikpgahefdbcgoiicnmpbeeomffld).

6. Browse to the URL you want to profile.

7. Open the extension by clicking on its icon in the browser toolbar and click the “Profile” button.

8. Once profiling is complete, you can click on the buttons in the toolbar to see the details of the profile.
